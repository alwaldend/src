package evaluation

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// boundedOutput caps a child event stream before it can exhaust controller memory.
type boundedOutput struct {
	// buffer owns the bounded bytes returned to the event decoder.
	buffer *bytes.Buffer
	// remaining is the maximum number of bytes still accepted.
	remaining int
}

// Write implements the subprocess output contract with an explicit byte limit.
func (output *boundedOutput) Write(data []byte) (int, error) {
	if len(data) > output.remaining {
		return 0, errors.New("ComfyUI event output exceeds 32 MiB")
	}
	count, err := output.buffer.Write(data)
	output.remaining -= count
	if err != nil {
		return count, fmt.Errorf("capture ComfyUI events: %w", err)
	}
	return count, nil
}

// evaluationCase owns one prompt, its supplied references and review rubric.
type evaluationCase struct {
	// ID names the case and its reproducibility record.
	ID string `json:"id"`
	// PromptFile is the case-relative positive prompt text.
	PromptFile string `json:"prompt_file"`
	// NegativePrompt is shared across all frames or one still image.
	NegativePrompt string `json:"negative_prompt"`
	// References lists the original local case inputs.
	References []reference `json:"references"`
	// Width selects the native latent and reference scaling width.
	Width int `json:"width"`
	// Height selects the native latent and reference scaling height.
	Height int `json:"height"`
	// FPS selects animation cadence when frames are configured.
	FPS float64 `json:"fps"`
	// Frames supplies one prompt for each chained img2img frame.
	Frames []string `json:"frame_prompts"`
	// Review defines semantic and artistic criteria without invented scores.
	Review []map[string]string `json:"review_criteria"`
	// EditRegion localizes edits in masked workflows using output-image pixels.
	EditRegion *editRegion `json:"edit_region,omitempty"`
}

// editRegion owns a case-specific feathered edit rectangle, independent of the model.
type editRegion struct {
	// X locates the rectangle horizontally in output pixels.
	X int `json:"x"`
	// Y locates the rectangle vertically.
	Y int `json:"y"`
	// Width is the horizontal extent of the editable region.
	Width int `json:"width"`
	// Height is the vertical extent of the editable region.
	Height int `json:"height"`
	// Feather softens the edit boundary without extending its rectangle.
	Feather int `json:"feather"`
}

// reference names one supplied image and its intended conditioning role.
type reference struct {
	// Path is a case-relative image filename.
	Path string `json:"path"`
	// Role explains what should be borrowed from that image.
	Role string `json:"role"`
}

// loadCase validates local case paths and decodes the maintained case definition.
func loadCase(options Options, job job) (evaluationCase, string, error) {
	var definition evaluationCase
	directory := filepath.Join(options.Cases, job.CaseID)
	if err := readJSON(filepath.Join(directory, "case.json"), &definition); err != nil {
		return definition, "", fmt.Errorf("load case %q: %w", job.CaseID, err)
	}
	if definition.ID != job.CaseID || definition.Width < 64 || definition.Height < 64 || definition.Width > 4096 || definition.Height > 4096 || definition.Width%16 != 0 || definition.Height%16 != 0 || filepath.Base(definition.PromptFile) != definition.PromptFile {
		return definition, "", errors.New("case has invalid identity, dimensions or prompt path")
	}
	data, err := os.ReadFile(filepath.Join(directory, definition.PromptFile))
	if err != nil {
		return definition, "", fmt.Errorf("read case prompt: %w", err)
	}
	if strings.TrimSpace(string(data)) == "" || len(data) > 65536 {
		return definition, "", errors.New("case prompt must be nonempty and at most 64 KiB")
	}
	for _, ref := range definition.References {
		if filepath.Base(ref.Path) != ref.Path || ref.Path == "." || ref.Path == ".." {
			return definition, "", errors.New("reference paths must be case-local filenames")
		}
		file, err := os.Open(filepath.Join(directory, ref.Path))
		if err != nil {
			return definition, "", fmt.Errorf("open reference %q: %w", ref.Path, err)
		}
		_, _, decodeErr := image.DecodeConfig(file)
		closeErr := file.Close()
		if err := errors.Join(decodeErr, closeErr); err != nil {
			return definition, "", fmt.Errorf("validate reference %q: %w", ref.Path, err)
		}
	}
	if len(definition.Frames) > 0 && (len(definition.Frames) < 2 || len(definition.Frames) > 120 || definition.FPS <= 0 || definition.FPS > 60) {
		return definition, "", errors.New("animation cases require 2–120 frame prompts and a cadence of 0–60 fps")
	}
	for index, frame := range definition.Frames {
		if strings.TrimSpace(frame) == "" || len(frame) > 65536 {
			return definition, "", fmt.Errorf("frame %d must be nonempty and at most 64 KiB", index+1)
		}
	}
	if definition.Width*definition.Height*max(1, len(definition.Frames)) > 64<<20 {
		return definition, "", errors.New("case exceeds the 64-million decoded-pixel budget")
	}
	if region := definition.EditRegion; region != nil && (region.X < 0 || region.Y < 0 || region.Width < 1 || region.Height < 1 || region.Width > definition.Width || region.Height > definition.Height || region.X > definition.Width-region.Width || region.Y > definition.Height-region.Height || region.Feather < 0 || region.Feather > min(region.Width, region.Height)/2) {
		return definition, "", errors.New("case edit region must fit inside the image with bounded feathering")
	}
	return definition, strings.TrimSpace(string(data)), nil
}

// validateWorkload checks native inputs and executable paths before rental allocation.
func validateWorkload(options Options, job job, definitionModel model) error {
	definition, prompt, err := loadCase(options, job)
	if err != nil {
		return fmt.Errorf("validate local case: %w", err)
	}
	var references []string
	for _, reference := range definition.References {
		references = append(references, reference.Path)
	}
	workflow, err := prepareWorkflow(options, job, definition, prompt, job.Seeds[0], references)
	if err != nil {
		return fmt.Errorf("validate workflow bindings: %w", err)
	}
	available := map[string]bool{}
	for _, name := range definitionModel.Downloads {
		available[name] = true
	}
	for _, value := range workflow["nodes"].([]any) {
		node := value.(map[string]any)
		kind, _ := node["type"].(string)
		count := 0
		switch kind {
		case "CheckpointLoaderSimple", "UNETLoader", "VAELoader", "CLIPLoader", "IPAdapterModelLoader", "CLIPVisionLoader":
			count = 1
		case "QuadrupleCLIPLoader":
			count = 4
		}
		if count > 0 {
			widgets, _ := node["widgets_values"].([]any)
			if len(widgets) < count {
				return errors.New("model loader is missing its declared files")
			}
			for _, value := range widgets[:count] {
				name, ok := value.(string)
				if !ok || !available[name] {
					return fmt.Errorf("loader %s selects a file absent from the model profile", kind)
				}
			}
		}
	}
	for _, path := range []string{options.Terraform, options.Bridge, options.CLI, options.CAFile} {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("inspect packaged tool %q: %w", path, err)
		}
		if info.IsDir() {
			return fmt.Errorf("packaged tool %q is a directory", path)
		}
	}
	return nil
}

// prepareWorkflow binds case inputs to the existing generated native UI graph.
func prepareWorkflow(options Options, job job, definition evaluationCase, prompt string, seed int64, references []string) (map[string]any, error) {
	var workflow map[string]any
	if err := readJSON(filepath.Join(options.Workflows, job.ModelProfile, job.Workflow+".json"), &workflow); err != nil {
		return nil, fmt.Errorf("read generated workflow: %w", err)
	}
	nodes, ok := workflow["nodes"].([]any)
	if !ok {
		return nil, errors.New("evaluation requires a generated native UI workflow")
	}
	positive, negative, loads, samplers, sequences := 0, 0, 0, 0, 0
	nativeSamplers, adapters, guides, gifs := 0, 0, 0, 0
	masks := 0
	for _, value := range nodes {
		node, ok := value.(map[string]any)
		if !ok {
			return nil, errors.New("workflow node is not an object")
		}
		widgets, ok := node["widgets_values"].([]any)
		if !ok {
			continue
		}
		typeName, _ := node["type"].(string)
		if controls := job.Sampling; controls != nil {
			if (typeName == "KSampler" || typeName == "ImagegenChainedImageToImageGIF" || typeName == "ImagegenQwenChainedImageToImageGIF") && len(widgets) >= 7 {
				if controls.Steps != nil {
					widgets[2] = *controls.Steps
				}
				if controls.CFG != nil {
					widgets[3] = *controls.CFG
				}
				if controls.Denoise != nil {
					widgets[6] = *controls.Denoise
				}
			}
			if typeName == "ClownsharKSampler_Beta" && len(widgets) >= 7 {
				if controls.Steps != nil {
					widgets[3] = *controls.Steps
				}
				if controls.CFG != nil {
					widgets[6] = *controls.CFG
				}
			}
			if typeName == "IPAdapterAdvanced" && len(widgets) >= 5 {
				if controls.AdapterWeight != nil {
					widgets[0] = *controls.AdapterWeight
				}
				if controls.AdapterEnd != nil {
					widgets[4] = *controls.AdapterEnd
				}
			}
			if typeName == "ClownGuide_Style_Beta" && len(widgets) >= 3 && controls.StyleWeight != nil {
				widgets[2] = *controls.StyleWeight
			}
		}
		switch typeName {
		case "SetLatentNoiseMask":
			masks++
		case "SolidMask":
			if len(widgets) < 3 || definition.EditRegion == nil {
				return nil, errors.New("masked workflow requires a case edit region")
			}
			if widgets[0] == float64(0) {
				widgets[1], widgets[2] = definition.Width, definition.Height
			} else {
				widgets[1], widgets[2] = definition.EditRegion.Width, definition.EditRegion.Height
			}
		case "FeatherMask":
			if len(widgets) != 4 || definition.EditRegion == nil {
				return nil, errors.New("feathered workflow requires a case edit region")
			}
			for index := range widgets {
				widgets[index] = definition.EditRegion.Feather
			}
		case "MaskComposite":
			if len(widgets) < 3 || definition.EditRegion == nil {
				return nil, errors.New("mask placement requires a case edit region")
			}
			widgets[0], widgets[1] = definition.EditRegion.X, definition.EditRegion.Y
		case "IPAdapterAdvanced":
			adapters++
		case "ClownGuide_Style_Beta":
			guides++
		case "CLIPTextEncode", "TextEncodeQwenImageEditPlus":
			if len(widgets) < 1 {
				return nil, errors.New("text encoder is missing its widget")
			}
			if positive == 0 {
				widgets[0] = prompt
				positive++
			} else {
				widgets[0] = definition.NegativePrompt
				negative++
			}
		case "ImagegenPromptSequence":
			if len(widgets) < 5 {
				return nil, errors.New("sequence node has incompatible widgets")
			}
			if len(definition.Frames) == 0 {
				return nil, errors.New("sequence workflow requires a frame-prompt case")
			}
			encodedFrames, err := json.Marshal(definition.Frames)
			if err != nil {
				return nil, fmt.Errorf("encode ordered frame prompts: %w", err)
			}
			widgets[0] = string(encodedFrames)
			widgets[1] = prompt
			widgets[2] = "JSON"
			widgets[3] = "evaluation"
			widgets[4] = definition.NegativePrompt
			sequences++
		case "LoadImage":
			if len(widgets) < 1 || len(references) == 0 {
				return nil, errors.New("reference workflow requires uploaded case images")
			}
			widgets[0] = references[loads%len(references)]
			loads++
		case "EmptyLatentImage", "EmptySD3LatentImage":
			if len(widgets) < 3 {
				return nil, errors.New("latent node has incompatible dimensions")
			}
			widgets[0] = definition.Width
			widgets[1] = definition.Height
			widgets[2] = 1
		case "ImageScale":
			if len(widgets) < 4 {
				return nil, errors.New("reference scaler has incompatible widgets")
			}
			// Adapter embeddings retain their deliberately fixed 512-square crop.
			if widgets[1] != float64(512) {
				widgets[1] = definition.Width
				widgets[2] = definition.Height
			}
		case "KSampler", "ImagegenChainedImageToImageGIF", "ImagegenQwenChainedImageToImageGIF":
			nativeSamplers++
			if len(widgets) < 7 {
				return nil, errors.New("sampler is missing a seed")
			}
			widgets[0] = seed
			widgets[1] = "fixed"
			samplers++
			if typeName == "ImagegenChainedImageToImageGIF" || typeName == "ImagegenQwenChainedImageToImageGIF" {
				gifs++
				if typeName == "ImagegenQwenChainedImageToImageGIF" {
					positive++
					negative++
				}
				if len(widgets) < 9 || len(definition.Frames) == 0 {
					return nil, errors.New("GIF workflow requires a frame case and compatible node")
				}
				widgets[7] = definition.FPS
				widgets[8] = "evaluation"
				if controls := job.Sampling; controls != nil && (controls.FrameSource != nil || controls.FrameSeedMode != nil) {
					if len(widgets) < 11 {
						return nil, errors.New("GIF frame overrides require the current node controls")
					}
					if controls.FrameSource != nil {
						widgets[9] = *controls.FrameSource
					}
					if controls.FrameSeedMode != nil {
						widgets[10] = *controls.FrameSeedMode
					}
				}
			}
		case "ClownsharKSampler_Beta":
			if len(widgets) < 9 {
				return nil, errors.New("style sampler has incompatible widgets")
			}
			widgets[7] = seed
			widgets[8] = "fixed"
			samplers++
		case "SaveImage":
			if len(widgets) > 0 {
				widgets[0] = "evaluation"
			}
		}
	}
	if positive != 1 || negative != 1 || samplers != 1 || (len(definition.Frames) > 0 && sequences != 1) {
		return nil, errors.New("workflow lacks a supported positive/negative sampler path for this case")
	}
	if len(definition.Frames) > 0 && gifs != 1 {
		return nil, errors.New("frame-prompt cases require a GIF output workflow")
	}
	if strings.HasSuffix(job.Workflow, "-masked") && masks != 1 {
		return nil, errors.New("masked workflow must connect a native latent noise mask")
	}
	if loads > 0 && len(definition.References) > loads {
		return nil, errors.New("case has more reference images than this workflow can condition on")
	}
	if controls := job.Sampling; controls != nil {
		if (controls.FrameSource != nil || controls.FrameSeedMode != nil) && gifs == 0 {
			return nil, errors.New("frame source and noise overrides require a GIF workflow")
		}
		if controls.Denoise != nil && nativeSamplers == 0 {
			return nil, errors.New("denoise override requires a native sampler")
		}
		if (controls.AdapterWeight != nil || controls.AdapterEnd != nil) && adapters == 0 {
			return nil, errors.New("adapter overrides require an adapter workflow")
		}
		if controls.StyleWeight != nil && guides == 0 {
			return nil, errors.New("style weight requires a style-guide workflow")
		}
	}
	return workflow, nil
}

// callCLI uses the existing authenticated bridge and retains its TLS/routing policy.
func callCLI(ctx context.Context, options Options, endpoint string, arguments ...string) (string, error) {
	argv := []string{"--cli", options.CLI, "--url", endpoint, "--ca-file", options.CAFile, "--"}
	argv = append(argv, arguments...)
	child := processCommand(ctx, options.Bridge, argv...)
	child.Env = workloadEnvironment()
	var output bytes.Buffer
	child.Stdout = &boundedOutput{buffer: &output, remaining: 32 << 20}
	child.Stderr = io.Discard
	if err := child.Run(); err != nil {
		code, _ := cliField(output.String(), "code")
		if len(code) > 64 {
			code = ""
		}
		for _, letter := range code {
			if !(letter >= 'a' && letter <= 'z' || letter >= '0' && letter <= '9' || letter == '_') {
				code = ""
				break
			}
		}
		if code != "" {
			return output.String(), fmt.Errorf("authenticated ComfyUI command failed (%s): %w", code, err)
		}
		return output.String(), fmt.Errorf("authenticated ComfyUI command failed: %w", err)
	}
	return output.String(), nil
}

// field recursively selects a named scalar from a CLI envelope or NDJSON event.
func field(value any, name string) string {
	switch value := value.(type) {
	case map[string]any:
		if result, ok := value[name].(string); ok && result != "" {
			return result
		}
		for _, child := range value {
			if result := field(child, name); result != "" {
				return result
			}
		}
	case []any:
		for _, child := range value {
			if result := field(child, name); result != "" {
				return result
			}
		}
	}
	return ""
}

// cliField reads either a full JSON envelope or a named NDJSON event value.
func cliField(output, name string) (string, error) {
	var value any
	if json.Unmarshal([]byte(output), &value) == nil {
		if result := field(value, name); result != "" {
			return result, nil
		}
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Buffer(make([]byte, 65536), 16<<20)
	for scanner.Scan() {
		if json.Unmarshal(scanner.Bytes(), &value) == nil {
			if result := field(value, name); result != "" {
				return result, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("decode CLI event stream: %w", err)
	}
	return "", fmt.Errorf("CLI result omitted %q", name)
}

// submitWorkload stages inputs and submits a seeded graph to the native queue.
func submitWorkload(ctx context.Context, options Options, job job, seed int64, endpoint, directory string) (queuedTrial, error) {
	trial := queuedTrial{job: job, seed: seed, directory: directory, started: time.Now()}
	definition, prompt, err := loadCase(options, job)
	if err != nil {
		return trial, fmt.Errorf("load execution case: %w", err)
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return trial, fmt.Errorf("create workload output: %w", err)
	}

	var references []string
	inputHashes := map[string]string{}
	for _, ref := range definition.References {
		path := filepath.Join(options.Cases, job.CaseID, ref.Path)
		data, err := os.ReadFile(path)
		if err != nil {
			return trial, fmt.Errorf("read input hash: %w", err)
		}
		sum := sha256.Sum256(data)
		inputHashes[ref.Path] = hex.EncodeToString(sum[:])
		output, err := callCLI(ctx, options, endpoint, "--json", "upload", path)
		if err != nil {
			return trial, fmt.Errorf("upload reference %q: %w", ref.Path, err)
		}
		name, err := cliField(output, "cloud_name")
		if err != nil {
			return trial, fmt.Errorf("read uploaded reference identity: %w", err)
		}
		references = append(references, name)
	}
	workflow, err := prepareWorkflow(options, job, definition, prompt, seed, references)
	if err != nil {
		return trial, fmt.Errorf("bind execution workflow: %w", err)
	}
	workflowPath := filepath.Join(directory, "workflow.json")
	if err := writeJSON(workflowPath, workflow); err != nil {
		return trial, fmt.Errorf("save effective workflow: %w", err)
	}
	generationStarted := time.Now()
	output, err := callCLI(ctx, options, endpoint, "run", "--workflow", workflowPath, "--json")
	if err != nil {
		return trial, fmt.Errorf("generate configured workload: %w", errors.Join(err, recordCLIFailure(directory, output)))
	}
	if err := os.WriteFile(filepath.Join(directory, "events.ndjson"), []byte(output), 0o600); err != nil {
		return trial, fmt.Errorf("save execution events: %w", err)
	}
	promptID, err := cliField(output, "prompt_id")
	if err != nil {
		return trial, fmt.Errorf("read submitted job identity: %w", err)
	}
	trial.definition = definition
	trial.promptID = promptID
	trial.hashes = inputHashes
	trial.generationStarted = generationStarted
	trial.events = output
	return trial, nil
}

// finalizeMedia verifies downloaded artifacts after the paid rental has been released.
func finalizeMedia(options Options, job job, directory string) error {
	definition, _, err := loadCase(options, job)
	if err != nil {
		return fmt.Errorf("load media criteria: %w", err)
	}
	media, err := inspectMedia(filepath.Join(directory, "media"), len(definition.Frames))
	if err != nil {
		return fmt.Errorf("verify retrieved media: %w", err)
	}
	for _, record := range media {
		if record["width"] != definition.Width || record["height"] != definition.Height {
			return fmt.Errorf("retrieved media dimensions do not match case %q", definition.ID)
		}
	}
	var report map[string]any
	path := filepath.Join(directory, "result.json")
	if err := readJSON(path, &report); err != nil {
		return fmt.Errorf("read retrieval report: %w", err)
	}
	report["media"] = media
	report["media_verified"] = true
	if err := writeJSON(path, report); err != nil {
		return fmt.Errorf("save media verification: %w", err)
	}
	return nil
}

// inspectMedia verifies decodability, dimensions and actual looping GIF frames.
func inspectMedia(directory string, expectedFrames int) ([]map[string]any, error) {
	var media []map[string]any
	err := filepath.WalkDir(directory, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk retrieved media: %w", walkErr)
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return errors.New("retrieved media must not be a symlink")
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect retrieved file: %w", err)
		}
		if info.Size() > 256<<20 {
			return errors.New("retrieved file exceeds 256 MiB")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read retrieved media: %w", err)
		}
		sum := sha256.Sum256(data)
		record := map[string]any{"file": entry.Name(), "sha256": hex.EncodeToString(sum[:]), "bytes": len(data)}
		if strings.EqualFold(filepath.Ext(path), ".gif") {
			animation, err := gif.DecodeAll(bytes.NewReader(data))
			if err != nil {
				return fmt.Errorf("decode generated GIF: %w", err)
			}
			if len(animation.Image) != expectedFrames || animation.LoopCount != 0 {
				return errors.New("GIF must preserve every configured frame and loop indefinitely")
			}
			unique := map[string]bool{}
			for index, frame := range animation.Image {
				unique[frameSignature(frame)] = true
				if animation.Delay[index] <= 0 {
					return errors.New("GIF contains a frame without a positive duration")
				}
			}
			record["frames"] = len(animation.Image)
			record["unique_frames"] = len(unique)
			record["pixel_change_detected"] = len(unique) > 1
			first, last := animation.Image[0], animation.Image[len(animation.Image)-1]
			if first.Bounds() == last.Bounds() && first.Bounds() == image.Rect(0, 0, animation.Config.Width, animation.Config.Height) {
				record["loop_pixel_difference"] = loopDifference(first, last)
			}
			record["loop_count"] = animation.LoopCount
			record["delays_centiseconds"] = animation.Delay
			record["width"] = animation.Config.Width
			record["height"] = animation.Config.Height
		} else {
			decoded, format, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				return fmt.Errorf("decode generated image: %w", err)
			}
			record["format"] = format
			record["width"] = decoded.Bounds().Dx()
			record["height"] = decoded.Bounds().Dy()
		}
		media = append(media, record)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("inspect output media: %w", err)
	}
	if len(media) == 0 {
		return nil, errors.New("completed workload downloaded no media")
	}
	if expectedFrames > 0 {
		found := false
		for _, record := range media {
			if record["frames"] != nil {
				found = true
			}
		}
		if !found {
			return nil, errors.New("animation workload downloaded no GIF")
		}
	}
	return media, nil
}

// frameSignature compares visible frame colors rather than palette index assignments.
func frameSignature(frame image.Image) string {
	bounds := frame.Bounds()
	pixels := make([]byte, bounds.Dx()*bounds.Dy()*4)
	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			red, green, blue, alpha := frame.At(x, y).RGBA()
			pixels[index] = byte(red >> 8)
			pixels[index+1] = byte(green >> 8)
			pixels[index+2] = byte(blue >> 8)
			pixels[index+3] = byte(alpha >> 8)
			index += 4
		}
	}
	sum := sha256.Sum256(pixels)
	return fmt.Sprintf("%v:%x", bounds, sum)
}

// loopDifference measures normalized boundary pixel change without judging natural motion.
func loopDifference(first, last image.Image) float64 {
	bounds := first.Bounds()
	total := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := first.At(x, y).RGBA()
			r2, g2, b2, _ := last.At(x, y).RGBA()
			total += math.Abs(float64(r1)-float64(r2)) + math.Abs(float64(g1)-float64(g2)) + math.Abs(float64(b1)-float64(b2))
		}
	}
	return total / (float64(bounds.Dx()*bounds.Dy()) * 3 * 65535)
}

// processCommand terminates the complete foreground process group on cancellation.
func processCommand(ctx context.Context, program string, args ...string) *exec.Cmd {
	child := exec.CommandContext(ctx, program, args...)
	child.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	child.WaitDelay = time.Second
	child.Cancel = func() error {
		if err := syscall.Kill(-child.Process.Pid, syscall.SIGKILL); err != nil {
			if errors.Is(err, syscall.ESRCH) {
				return fmt.Errorf("process group already ended: %w", os.ErrProcessDone)
			}
			return fmt.Errorf("terminate workload process group: %w", err)
		}
		return nil
	}
	return child
}

// queuedTrial retains one submitted workflow until its own results are retrieved.
type queuedTrial struct {
	// completion receives independently observed termination, unaffected by preceding downloads.
	completion <-chan trialCompletion
	// deadline bounds submission, native queue residence, execution and retrieval together.
	deadline time.Time
	// job preserves the expanded matrix coordinate and sampling settings.
	job job
	// seed selects this independent fixed trial.
	seed int64
	// definition owns its media criteria and review rubric.
	definition evaluationCase
	// directory owns this trial's workflow, events and result media.
	directory string
	// promptID identifies this native queue entry.
	promptID string
	// hashes records exact reference bytes used by the submission.
	hashes map[string]string
	// started includes input upload, queueing, generation and retrieval.
	started time.Time
	// generationStarted includes native queue residence before inference.
	generationStarted time.Time
	// events retains the initial CLI submission envelope.
	events string
}

// trialCompletion retains a native completion observation independently of retrieval order.
type trialCompletion struct {
	// observedAt includes CLI/polling latency but excludes preceding download delays.
	observedAt time.Time
	// events is the validated compact terminal summary.
	events string
	// err preserves completion or diagnostic failures.
	err error
}

// observeTrial watches one queued job with its original absolute deadline.
func observeTrial(ctx context.Context, options Options, trial queuedTrial, endpoint string, poll time.Duration) trialCompletion {
	timeout := 1800.0
	if deadline, ok := ctx.Deadline(); ok {
		timeout = max(0.001, time.Until(deadline).Seconds())
	}
	output, err := callCLI(ctx, options, endpoint, "--json", "jobs", "wait", trial.promptID, "--poll-interval", strconv.FormatFloat(poll.Seconds(), 'f', 3, 64), "--timeout", strconv.FormatFloat(timeout, 'f', 3, 64))
	if err != nil {
		// The wait summary omits native node exceptions; capture status before cleanup.
		diagnostic, diagnosticCancel := context.WithTimeout(ctx, 10*time.Second)
		details, _ := callCLI(diagnostic, options, endpoint, "--json", "jobs", "status", trial.promptID)
		diagnosticCancel()
		return trialCompletion{err: fmt.Errorf("wait for queued workflow: %w", errors.Join(err, recordCLIFailure(trial.directory, output+"\n"+details)))}
	}
	var receipt struct {
		// Data is the pinned CLI's terminal queue summary.
		Data struct {
			// Jobs contains exactly the one requested native queue entry.
			Jobs []struct {
				// PromptID identifies the requested native queue entry.
				PromptID string `json:"prompt_id"`
				// Status must confirm successful completion before retrieval.
				Status string `json:"status"`
				// OK is the CLI's success acknowledgment for that entry.
				OK bool `json:"ok"`
			} `json:"jobs"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(output), &receipt); err != nil {
		return trialCompletion{err: fmt.Errorf("decode terminal queue summary: %w", err)}
	}
	if len(receipt.Data.Jobs) != 1 || receipt.Data.Jobs[0].PromptID != trial.promptID || receipt.Data.Jobs[0].Status != "completed" || !receipt.Data.Jobs[0].OK {
		return trialCompletion{err: errors.New("native queue did not confirm successful workflow completion")}
	}
	observedAt := time.Now()
	var compact bytes.Buffer
	if err := json.Compact(&compact, []byte(output)); err != nil {
		return trialCompletion{err: fmt.Errorf("compact terminal queue event: %w", err)}
	}
	return trialCompletion{observedAt: observedAt, events: compact.String()}
}

// retrieveTrial downloads an observed completion within the original trial deadline.
func retrieveTrial(ctx context.Context, options Options, trial queuedTrial, completion trialCompletion, endpoint string) error {
	if completion.err != nil {
		return completion.err
	}
	if err := os.WriteFile(filepath.Join(trial.directory, "events.ndjson"), []byte(strings.TrimSpace(trial.events)+"\n"+completion.events+"\n"), 0o600); err != nil {
		return fmt.Errorf("save queue completion events: %w", err)
	}
	generationSeconds := completion.observedAt.Sub(trial.generationStarted).Seconds()
	downloadStarted := time.Now()
	if output, err := callCLI(ctx, options, endpoint, "--json", "download", trial.promptID, "-o", filepath.Join(trial.directory, "media")); err != nil {
		return fmt.Errorf("retrieve workload results: %w", errors.Join(err, recordCLIFailure(trial.directory, output)))
	}
	report := map[string]any{"case": trial.definition, "job": trial.job, "seed": trial.seed, "prompt_id": trial.promptID, "reference_sha256": trial.hashes, "deadline_at": trial.deadline.UTC(), "completion_observed_at": completion.observedAt.UTC(), "timing_source": "independent_cli_observation", "queue_and_generation_seconds": generationSeconds, "retrieval_seconds": time.Since(downloadStarted).Seconds(), "total_seconds": time.Since(trial.started).Seconds(), "download_complete": true, "media_verified": false, "review_scores": nil}
	if err := writeJSON(filepath.Join(trial.directory, "result.json"), report); err != nil {
		return fmt.Errorf("save workload results: %w", err)
	}
	return nil
}

// recordCLIFailure preserves bounded diagnostics locally after masking injected credentials.
func recordCLIFailure(directory, output string) error {
	if output == "" {
		return nil
	}
	for _, name := range []string{"IMAGEGEN_WEB_PASSWORD", "VAST_API_KEY", "HF_TOKEN", "VAULT_TOKEN"} {
		if value := os.Getenv(name); value != "" {
			output = strings.ReplaceAll(output, value, "[redacted]")
		}
	}
	if err := os.WriteFile(filepath.Join(directory, "cli-failure.ndjson"), []byte(output), 0o600); err != nil {
		return fmt.Errorf("preserve local CLI failure diagnostics: %w", err)
	}
	return nil
}

// queueTrials keeps a bounded lookahead of native work submitted to one rental.
func queueTrials(ctx context.Context, options Options, jobs []job, endpoint, directory string, timeout, poll time.Duration, trials chan<- queuedTrial, observers *sync.WaitGroup) error {
	for _, job := range jobs {
		for _, seed := range job.Seeds {
			if options.logger != nil {
				options.logger.Printf("workload %s: queueing seed %d", job.ID, seed)
			}
			submission, submissionCancel := context.WithTimeout(ctx, timeout)
			trial, err := submitWorkload(submission, options, job, seed, endpoint, filepath.Join(directory, job.ID, strconv.FormatInt(seed, 10)))
			trial.deadline, _ = submission.Deadline()
			submissionCancel()
			if err != nil {
				return fmt.Errorf("queue workload %q seed %d: %w", job.ID, seed, err)
			}
			completion := make(chan trialCompletion, 1)
			trial.completion = completion
			observers.Add(1)
			go func(trial queuedTrial, result chan<- trialCompletion) {
				defer observers.Done()
				observation, observationCancel := context.WithDeadline(ctx, trial.deadline)
				defer observationCancel()
				result <- observeTrial(observation, options, trial, endpoint, poll)
			}(trial, completion)
			select {
			case trials <- trial:
			case <-ctx.Done():
				return fmt.Errorf("cancel native queue submission: %w", ctx.Err())
			}
		}
	}
	return nil
}

// runQueuedWorkloads overlaps native queue submission with per-trial result retrieval.
func runQueuedWorkloads(ctx context.Context, options Options, jobs []job, endpoint, directory string, timeout, poll time.Duration) error {
	ctx, cancel := context.WithCancel(ctx)
	var observers sync.WaitGroup
	defer func() { cancel(); observers.Wait() }()
	trials := make(chan queuedTrial, 2)
	submitted := make(chan error, 1)
	go func() {
		err := queueTrials(ctx, options, jobs, endpoint, directory, timeout, poll, trials, &observers)
		close(trials)
		submitted <- err
	}()
	var resultErr error
	for trial := range trials {
		if resultErr != nil {
			continue
		}
		completion := <-trial.completion
		retrieval, retrievalCancel := context.WithDeadline(ctx, trial.deadline)
		err := retrieveTrial(retrieval, options, trial, completion, endpoint)
		retrievalCancel()
		if err != nil {
			resultErr = fmt.Errorf("collect workload %q seed %d: %w", trial.job.ID, trial.seed, err)
			cancel()
			continue
		}
		if options.logger != nil {
			options.logger.Printf("workload %s: seed %d downloaded", trial.job.ID, trial.seed)
		}
	}
	if err := errors.Join(resultErr, <-submitted); err != nil {
		return fmt.Errorf("execute rental queue: %w", err)
	}
	return nil
}
