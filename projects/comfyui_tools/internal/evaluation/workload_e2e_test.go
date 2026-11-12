package evaluation

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"golang.org/x/net/websocket"
)

// TestWorkloadThroughCLI exercises real upload, native conversion, execution and retrieval.
func TestWorkloadThroughCLI(t *testing.T) {
	for _, mode := range []string{"image", "animation", "queue", "qwen_animation", "deadline", "empty_frame", "execution_error"} {
		t.Run(mode, func(t *testing.T) { exerciseCLIWorkload(t, mode) })
	}
}

// exerciseCLIWorkload uses real transport for either a still image or chained GIF.
func exerciseCLIWorkload(t *testing.T, mode string) {
	animation := mode == "animation" || mode == "qwen_animation" || mode == "empty_frame"
	queue, qwen, deadline := mode == "queue", mode == "qwen_animation", mode == "deadline"
	directory := t.TempDir()
	var imageBytes bytes.Buffer
	canvas := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			canvas.Set(x, y, color.RGBA{R: uint8(x * 4), B: uint8(y * 4), A: 255})
		}
	}
	if err := png.Encode(&imageBytes, canvas); err != nil {
		t.Fatalf("encode public image fixture: %v", err)
	}
	outputName, outputNode, outputKind := "fixture.png", "9", "images"
	outputBytes := imageBytes.Bytes()
	definition := evaluationCase{ID: "fixture", PromptFile: "prompt.txt", NegativePrompt: "letters", References: []reference{{Path: "reference.png", Role: "fixture"}}, Width: 64, Height: 64}
	workflow := fixtureWorkflow()
	if animation {
		outputName, outputNode, outputKind = "fixture.gif", "8", "gifs"
		first := image.NewPaletted(image.Rect(0, 0, 64, 64), color.Palette{color.RGBA{B: 255, A: 255}})
		second := image.NewPaletted(image.Rect(0, 0, 64, 64), color.Palette{color.RGBA{G: 255, A: 255}})
		var encoded bytes.Buffer
		if err := gif.EncodeAll(&encoded, &gif.GIF{Image: []*image.Paletted{first, second, first}, Delay: []int{17, 17, 17}, LoopCount: 0}); err != nil {
			t.Fatalf("encode controlled GIF: %v", err)
		}
		outputBytes = encoded.Bytes()
		definition.Frames = []string{"The cube tilts left.\nThe shadow stays still.", "The cube tilts right.", "The cube returns to the initial pose."}
		if mode == "empty_frame" {
			definition.Frames[1] = " \n "
		}
		definition.FPS = 6
		workflow = fixtureAnimationWorkflow()
		if qwen {
			workflow = fixtureQwenAnimationWorkflow()
		}
	}
	outputData := map[string]any{outputKind: []map[string]any{{"filename": outputName, "subfolder": "", "type": "output"}}}
	queued := make(chan struct{})
	var once sync.Once
	var guard sync.Mutex
	var submitted map[string]any
	submissions := 0
	uploads := 0
	firstDownloadSubmissions := 0
	var firstDownloadFinishedAt time.Time
	info := fixtureNodeInfo()
	history := map[string]any{"fixture-job": map[string]any{"outputs": map[string]any{outputNode: outputData}, "status": map[string]any{"status_str": "success", "completed": true, "messages": []any{}}}}

	if mode == "execution_error" {
		history["fixture-job"] = map[string]any{"outputs": map[string]any{}, "status": map[string]any{"status_str": "error", "completed": false, "messages": []any{[]any{"execution_error", map[string]any{"node_id": "7", "node_type": "KSampler", "exception_type": "RuntimeError", "exception_message": "fixture execution failed fixture-web-password"}}}}}
	}

	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "fixture" || password != "fixture-web-password" {
			http.Error(w, "authentication required", 401)
			return
		}
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, "ComfyUI fixture")
		case "/history":
			fmt.Fprint(w, `{}`)
		case "/upload/image":
			guard.Lock()
			uploads++
			guard.Unlock()
			if deadline {
				time.Sleep(1700 * time.Millisecond)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"name":"uploaded-reference.png","subfolder":"","type":"input"}`)
		case "/object_info":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(info); err != nil {
				return
			}
		case "/prompt":
			if r.Method == http.MethodGet {
				fmt.Fprint(w, `{"exec_info":{"queue_remaining":0}}`)
				return
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid request", 400)
				return
			}
			guard.Lock()
			submitted, _ = body["prompt"].(map[string]any)
			guard.Unlock()
			guard.Lock()
			submissions++
			id := fmt.Sprintf("fixture-job-%d", submissions)
			guard.Unlock()
			fmt.Fprintf(w, `{"prompt_id":%q,"number":1,"node_errors":{}}`, id)
			once.Do(func() { close(queued) })
		case "/history/fixture-job-1", "/history/fixture-job-2":
			guard.Lock()
			ready := !queue || submissions >= 2
			guard.Unlock()
			if !ready {
				fmt.Fprint(w, `{}`)
				return
			}
			response := map[string]any{strings.TrimPrefix(r.URL.Path, "/history/"): history["fixture-job"]}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				return
			}
		case "/queue":
			guard.Lock()
			pending := queue && submissions < 2
			guard.Unlock()
			if pending {
				fmt.Fprint(w, `{"queue_running":[[1,"fixture-job-1",{},{}]],"queue_pending":[]}`)
				return
			}
			fmt.Fprint(w, `{"queue_running":[],"queue_pending":[]}`)
		case "/view":
			if deadline {
				time.Sleep(1700 * time.Millisecond)
			}
			guard.Lock()
			first := firstDownloadSubmissions == 0
			if first {
				firstDownloadSubmissions = submissions
			}
			guard.Unlock()
			if queue && first {
				time.Sleep(2 * time.Second)
			}
			guard.Lock()
			if first {
				firstDownloadFinishedAt = time.Now()
			}
			guard.Unlock()
			w.Header().Set("Content-Type", "application/octet-stream")
			if _, err := w.Write(outputBytes); err != nil {
				return
			}
		case "/system_stats":
			fmt.Fprint(w, `{"system":{"os":"fixture","comfyui_version":"fixture"},"devices":[]}`)
		case "/ws":
			websocket.Handler(func(connection *websocket.Conn) {
				defer connection.Close()
				select {
				case <-queued:
				case <-time.After(5 * time.Second):
					return
				}
				executed, err := json.Marshal(map[string]any{"type": "executed", "data": map[string]any{"node": outputNode, "prompt_id": "fixture-job", "output": outputData}})
				if err != nil {
					return
				}
				messages := []string{`{"type":"execution_start","data":{"prompt_id":"fixture-job"}}`, string(executed), `{"type":"executing","data":{"node":null,"prompt_id":"fixture-job"}}`}

				for _, message := range messages {
					if err := websocket.Message.Send(connection, message); err != nil {
						return
					}
				}
			}).ServeHTTP(w, r)
		default:
			http.Error(w, "unsupported fixture endpoint "+r.URL.Path, 404)
		}
	})
	server := httptest.NewTLSServer(api)
	defer server.Close()
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	ca := filepath.Join(directory, "ca.pem")
	if err := os.WriteFile(ca, certificate, 0o600); err != nil {
		t.Fatalf("write fixture CA: %v", err)
	}
	bridge, err := runfiles.Rlocation(os.Getenv("BRIDGE_BINARY"))
	if err != nil {
		t.Fatalf("resolve bridge: %v", err)
	}
	cli, err := runfiles.Rlocation(os.Getenv("CLI_BINARY"))
	if err != nil {
		t.Fatalf("resolve CLI: %v", err)
	}
	caseDirectory := filepath.Join(directory, "cases", "fixture")
	if err := os.MkdirAll(caseDirectory, 0o700); err != nil {
		t.Fatalf("create public case: %v", err)
	}
	if err := os.WriteFile(filepath.Join(caseDirectory, "reference.png"), imageBytes.Bytes(), 0o600); err != nil {
		t.Fatalf("write fixture reference: %v", err)
	}
	if err := os.WriteFile(filepath.Join(caseDirectory, "prompt.txt"), []byte("A blue cube on a table."), 0o600); err != nil {
		t.Fatalf("write fixture prompt: %v", err)
	}
	if err := writeJSON(filepath.Join(caseDirectory, "case.json"), definition); err != nil {
		t.Fatalf("write fixture case: %v", err)
	}
	workflowDirectory := filepath.Join(directory, "workflows", "fixture")
	if err := writeJSON(filepath.Join(workflowDirectory, "fixture.json"), workflow); err != nil {
		t.Fatalf("write native fixture graph: %v", err)
	}
	t.Setenv("IMAGEGEN_WEB_USERNAME", "fixture")
	t.Setenv("IMAGEGEN_WEB_PASSWORD", "fixture-web-password")
	options := Options{Cases: filepath.Join(directory, "cases"), Workflows: filepath.Join(directory, "workflows"), Bridge: bridge, CLI: cli, CAFile: ca}
	steps, cfg, denoise := 2, 2.0, 0.4
	selection := job{ID: "fixture", CaseID: "fixture", ModelProfile: "fixture", Workflow: "fixture", Seeds: []int64{42}, Sampling: &samplingOverrides{Steps: &steps, CFG: &cfg, Denoise: &denoise}}
	if animation {
		source, mode := "original", "fixed"
		selection.Sampling.FrameSource = &source
		selection.Sampling.FrameSeedMode = &mode
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output := filepath.Join(directory, "results")
	if queue {
		selection.Seeds = []int64{41, 42}
	}
	budget := 20 * time.Second
	if deadline {
		budget = 3 * time.Second
	}
	err = runQueuedWorkloads(ctx, options, []job{selection}, server.URL, output, budget, 20*time.Millisecond)
	if mode == "execution_error" {
		diagnostic, readErr := os.ReadFile(filepath.Join(output, selection.ID, "42", "cli-failure.ndjson"))
		if err == nil || readErr != nil || !bytes.Contains(diagnostic, []byte("fixture execution failed")) || bytes.Contains(diagnostic, []byte("fixture-web-password")) {
			t.Fatalf("failed workload lost its redacted native node exception: run=%v read=%v", err, readErr)
		}
		if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
			if err := os.WriteFile(filepath.Join(artifact, "cli-execution-error.ndjson"), diagnostic, 0o600); err != nil {
				t.Fatalf("preserve repeatable error diagnostic: %v", err)
			}
		}
		return
	}
	if deadline || mode == "empty_frame" {
		guard.Lock()
		defer guard.Unlock()
		if err == nil || (mode == "empty_frame" && (uploads != 0 || submissions != 0)) {
			t.Fatalf("invalid frame or shared trial budget was accepted: %v; uploads=%d submissions=%d", err, uploads, submissions)
		}
		if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
			if err := writeJSON(filepath.Join(artifact, "cli-"+mode+".json"), map[string]any{"rejected": true, "uploads": uploads, "submissions": submissions, "live_operations": 0}); err != nil {
				t.Fatalf("save failure receipt: %v", err)
			}
		}
		return
	}
	if err != nil {
		t.Fatalf("execute real CLI workload: %v", err)
	}
	output = filepath.Join(output, selection.ID, "42")
	if err := finalizeMedia(options, selection, output); err != nil {
		t.Fatalf("inspect downloaded output: %v", err)
	}
	guard.Lock()
	defer guard.Unlock()
	if queue && firstDownloadSubmissions < 2 {
		t.Fatal("workflows were not queued before awaiting the first result")
	}
	expectedInputs := map[string]map[string]any{"2": {"text": "A blue cube on a table."}, "3": {"text": "letters"}, "4": {"image": "uploaded-reference.png"}, "7": {"seed": float64(42), "steps": float64(2), "cfg": float64(2), "denoise": 0.4}}
	if animation {
		frameJSON, err := json.Marshal(definition.Frames)
		if err != nil {
			t.Fatalf("encode expected ordered frame prompts: %v", err)
		}
		expectedInputs = map[string]map[string]any{"7": {"prompts": string(frameJSON), "separator": "JSON", "context": "A blue cube on a table.", "negative_prompt": "letters"}, "8": {"seed": float64(42), "steps": float64(2), "cfg": float64(2), "denoise": 0.4, "fps": float64(6), "frame_source": "original", "frame_seed_mode": "fixed"}}
	}
	if qwen {
		inputs := submitted["8"].(map[string]any)["inputs"].(map[string]any)
		for _, name := range []string{"clip", "reference_image", "prompts", "negative_prompt", "vae"} {
			if len(inputs[name].([]any)) != 2 {
				t.Fatalf("Qwen GIF lost its live %s input", name)
			}
		}
		if submitted["8"].(map[string]any)["class_type"] != "ImagegenQwenChainedImageToImageGIF" {
			t.Fatal("Qwen GIF still uses precomputed edit conditionings")
		}
	}
	for id, expected := range expectedInputs {
		node, ok := submitted[id].(map[string]any)
		if !ok {
			t.Fatalf("submitted native workflow omitted node %s", id)
		}
		inputs, ok := node["inputs"].(map[string]any)
		if !ok {
			t.Fatalf("submitted node %s omitted its inputs", id)
		}
		for name, value := range expected {
			if inputs[name] != value {
				t.Fatalf("submitted node %s field %s differs from configured case", id, name)
			}
		}
	}
	data, err := os.ReadFile(filepath.Join(output, "result.json"))
	if err != nil {
		t.Fatalf("read workload report: %v", err)
	}
	if animation {
		var report struct {
			// Media records the verified downloaded frame metadata.
			Media []map[string]any `json:"media"`
		}
		if err := json.Unmarshal(data, &report); err != nil {
			t.Fatalf("decode GIF receipt: %v", err)
		}
		if len(report.Media) != 1 || report.Media[0]["unique_frames"] != float64(2) || report.Media[0]["loop_pixel_difference"] != float64(0) {
			t.Fatal("GIF metrics failed to distinguish palettes or preserve the loop boundary")
		}
	}
	if queue {
		var timing struct {
			// ObservedAt must come from the independent completion observer, before the preceding download ends.
			ObservedAt time.Time `json:"completion_observed_at"`
		}
		if err := json.Unmarshal(data, &timing); err != nil {
			t.Fatalf("decode completion observation: %v", err)
		}
		if timing.ObservedAt.IsZero() || !timing.ObservedAt.Before(firstDownloadFinishedAt) {
			t.Fatal("second job completion was measured only after the first download")
		}
	}
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.WriteFile(filepath.Join(artifact, "cli-workload-"+mode+".json"), data, 0o600); err != nil {
			t.Fatalf("preserve repeatable workload artifact: %v", err)
		}
	}
}

// fixtureNodeInfo supplies the native stock-node contracts consumed by comfy-cli.
func fixtureNodeInfo() map[string]any {
	definitions := map[string]map[string]any{
		"CheckpointLoaderSimple": {"ckpt_name": []any{[]string{"fixture.safetensors"}}},
		"CLIPTextEncode":         {"text": []any{"STRING", map[string]any{"multiline": true}}, "clip": []any{"CLIP"}},
		"LoadImage":              {"image": []any{[]string{"uploaded-reference.png"}, map[string]any{"image_upload": true}}},
		"ImageScale":             {"image": []any{"IMAGE"}, "upscale_method": []any{[]string{"lanczos"}}, "width": []any{"INT"}, "height": []any{"INT"}, "crop": []any{[]string{"center"}}},
		"VAEEncode":              {"pixels": []any{"IMAGE"}, "vae": []any{"VAE"}},
		"KSampler":               {"model": []any{"MODEL"}, "seed": []any{"INT", map[string]any{"control_after_generate": true}}, "steps": []any{"INT"}, "cfg": []any{"FLOAT"}, "sampler_name": []any{[]string{"euler"}}, "scheduler": []any{[]string{"normal"}}, "positive": []any{"CONDITIONING"}, "negative": []any{"CONDITIONING"}, "latent_image": []any{"LATENT"}, "denoise": []any{"FLOAT"}},
		"VAEDecode":              {"samples": []any{"LATENT"}, "vae": []any{"VAE"}},
		"SaveImage":              {"images": []any{"IMAGE"}, "filename_prefix": []any{"STRING"}},
	}
	orders := map[string][]string{"CheckpointLoaderSimple": {"ckpt_name"}, "CLIPTextEncode": {"text", "clip"}, "LoadImage": {"image"}, "ImageScale": {"image", "upscale_method", "width", "height", "crop"}, "VAEEncode": {"pixels", "vae"}, "KSampler": {"model", "seed", "steps", "cfg", "sampler_name", "scheduler", "positive", "negative", "latent_image", "denoise"}, "VAEDecode": {"samples", "vae"}, "SaveImage": {"images", "filename_prefix"}}
	outputs := map[string][]string{"CheckpointLoaderSimple": {"MODEL", "CLIP", "VAE"}, "CLIPTextEncode": {"CONDITIONING"}, "LoadImage": {"IMAGE", "MASK"}, "ImageScale": {"IMAGE"}, "VAEEncode": {"LATENT"}, "KSampler": {"LATENT"}, "VAEDecode": {"IMAGE"}, "SaveImage": {}}
	result := map[string]any{}
	for name, required := range definitions {
		result[name] = map[string]any{"input": map[string]any{"required": required}, "input_order": map[string]any{"required": orders[name]}, "output": outputs[name], "output_name": outputs[name], "output_node": name == "SaveImage", "name": name}
	}
	sequence := map[string]any{"input": map[string]any{"required": map[string]any{"prompts": []any{"STRING"}, "context": []any{"STRING"}, "separator": []any{[]string{"line", "blank line", "JSON"}}, "filename_prefix": []any{"STRING"}}, "optional": map[string]any{"negative_prompt": []any{"STRING"}, "model": []any{"MODEL"}}}, "input_order": map[string]any{"required": []string{"prompts", "context", "separator", "filename_prefix"}, "optional": []string{"negative_prompt", "model"}}, "output": []string{"STRING", "STRING", "STRING", "MODEL"}, "output_node": false}
	result["ImagegenPromptSequence"] = sequence
	gifInputs := map[string]any{"model": []any{"MODEL"}, "positive": []any{"CONDITIONING"}, "negative": []any{"CONDITIONING"}, "latent_image": []any{"LATENT"}, "vae": []any{"VAE"}, "seed": []any{"INT", map[string]any{"control_after_generate": true}}, "steps": []any{"INT"}, "cfg": []any{"FLOAT"}, "sampler_name": []any{[]string{"euler"}}, "scheduler": []any{[]string{"normal"}}, "denoise": []any{"FLOAT"}, "fps": []any{"FLOAT"}, "filename_prefix": []any{"STRING"}}
	result["ImagegenChainedImageToImageGIF"] = map[string]any{"input": map[string]any{"required": gifInputs}, "input_order": map[string]any{"required": []string{"model", "positive", "negative", "latent_image", "vae", "seed", "steps", "cfg", "sampler_name", "scheduler", "denoise", "fps", "filename_prefix"}}, "output": []string{"IMAGE"}, "output_node": true}
	gifDefinition := result["ImagegenChainedImageToImageGIF"].(map[string]any)
	gifDefinition["input"].(map[string]any)["optional"] = map[string]any{"frame_source": []any{[]string{"previous", "original"}}, "frame_seed_mode": []any{[]string{"increment", "fixed"}}}
	gifDefinition["input_order"].(map[string]any)["optional"] = []string{"frame_source", "frame_seed_mode"}
	result["TextEncodeQwenImageEditPlus"] = map[string]any{"input": map[string]any{"required": map[string]any{"clip": []any{"CLIP"}, "prompt": []any{"STRING"}, "vae": []any{"VAE"}}, "optional": map[string]any{"image1": []any{"IMAGE"}}}, "input_order": map[string]any{"required": []string{"clip", "prompt", "vae"}, "optional": []string{"image1"}}, "output": []string{"CONDITIONING"}, "output_node": false}
	qwenInputs := map[string]any{}
	for key, value := range gifInputs {
		if key != "positive" && key != "negative" {
			qwenInputs[key] = value
		}
	}
	qwenInputs["prompts"] = []any{"STRING", map[string]any{"forceInput": true}}
	qwenInputs["negative_prompt"] = []any{"STRING", map[string]any{"forceInput": true}}
	qwenInputs["clip"] = []any{"CLIP"}
	qwenInputs["reference_image"] = []any{"IMAGE"}
	result["ImagegenQwenChainedImageToImageGIF"] = map[string]any{"input": map[string]any{"required": qwenInputs, "optional": gifDefinition["input"].(map[string]any)["optional"]}, "input_order": map[string]any{"required": []string{"model", "prompts", "negative_prompt", "latent_image", "vae", "seed", "steps", "cfg", "sampler_name", "scheduler", "denoise", "fps", "filename_prefix", "clip", "reference_image"}, "optional": []string{"frame_source", "frame_seed_mode"}}, "output": []string{"IMAGE"}, "output_node": true}
	return result
}

// fixtureWorkflow is one connected native img2img graph with public dummy weights.
func fixtureWorkflow() map[string]any {
	nodes := []map[string]any{
		{"id": 1, "type": "CheckpointLoaderSimple", "widgets_values": []any{"fixture.safetensors"}},
		{"id": 2, "type": "CLIPTextEncode", "widgets_values": []any{"positive"}},
		{"id": 3, "type": "CLIPTextEncode", "widgets_values": []any{"negative"}},
		{"id": 4, "type": "LoadImage", "widgets_values": []any{"reference.png", "image"}},
		{"id": 5, "type": "ImageScale", "widgets_values": []any{"lanczos", 64, 64, "center"}},
		{"id": 6, "type": "VAEEncode", "widgets_values": []any{}},
		{"id": 7, "type": "KSampler", "widgets_values": []any{0, "randomize", 1, 1, "euler", "normal", 0.7}},
		{"id": 8, "type": "VAEDecode", "widgets_values": []any{}},
		{"id": 9, "type": "SaveImage", "widgets_values": []any{"fixture"}},
	}
	links := [][]any{{1, 1, 1, 2, 0, "CLIP"}, {2, 1, 1, 3, 0, "CLIP"}, {3, 1, 0, 7, 0, "MODEL"}, {4, 2, 0, 7, 1, "CONDITIONING"}, {5, 3, 0, 7, 2, "CONDITIONING"}, {6, 4, 0, 5, 0, "IMAGE"}, {7, 5, 0, 6, 0, "IMAGE"}, {8, 1, 2, 6, 1, "VAE"}, {9, 6, 0, 7, 3, "LATENT"}, {10, 7, 0, 8, 0, "LATENT"}, {11, 1, 2, 8, 1, "VAE"}, {12, 8, 0, 9, 0, "IMAGE"}}
	inputs := map[int][]map[string]any{2: {{"name": "clip", "type": "CLIP", "link": 1}}, 3: {{"name": "clip", "type": "CLIP", "link": 2}}, 5: {{"name": "image", "type": "IMAGE", "link": 6}}, 6: {{"name": "pixels", "type": "IMAGE", "link": 7}, {"name": "vae", "type": "VAE", "link": 8}}, 7: {{"name": "model", "type": "MODEL", "link": 3}, {"name": "positive", "type": "CONDITIONING", "link": 4}, {"name": "negative", "type": "CONDITIONING", "link": 5}, {"name": "latent_image", "type": "LATENT", "link": 9}}, 8: {{"name": "samples", "type": "LATENT", "link": 10}, {"name": "vae", "type": "VAE", "link": 11}}, 9: {{"name": "images", "type": "IMAGE", "link": 12}}}
	for _, node := range nodes {
		node["inputs"] = inputs[node["id"].(int)]
	}
	return map[string]any{"nodes": nodes, "links": links, "version": 0.4}
}

// fixtureAnimationWorkflow connects frame prompts through the native chained output node.
func fixtureAnimationWorkflow() map[string]any {
	workflow := fixtureWorkflow()
	nodes := workflow["nodes"].([]map[string]any)[:6]
	nodes[1]["inputs"] = []map[string]any{{"name": "clip", "type": "CLIP", "link": 1}, {"name": "text", "type": "STRING", "link": 12, "widget": map[string]string{"name": "text"}}}
	nodes[2]["inputs"] = []map[string]any{{"name": "clip", "type": "CLIP", "link": 2}, {"name": "text", "type": "STRING", "link": 13, "widget": map[string]string{"name": "text"}}}
	nodes = append(nodes, map[string]any{"id": 7, "type": "ImagegenPromptSequence", "widgets_values": []any{"frames", "context", "line", "evaluation", "negative"}, "inputs": []map[string]any{{"name": "model", "type": "MODEL", "link": 3}}})
	nodes = append(nodes, map[string]any{"id": 8, "type": "ImagegenChainedImageToImageGIF", "widgets_values": []any{0, "randomize", 1, 1, "euler", "normal", 0.45, 6, "evaluation", "previous", "increment"}, "inputs": []map[string]any{{"name": "model", "type": "MODEL", "link": 4}, {"name": "positive", "type": "CONDITIONING", "link": 5}, {"name": "negative", "type": "CONDITIONING", "link": 9}, {"name": "latent_image", "type": "LATENT", "link": 10}, {"name": "vae", "type": "VAE", "link": 11}, {"name": "filename_prefix", "type": "STRING", "link": 14, "widget": map[string]string{"name": "filename_prefix"}}}})
	workflow["nodes"] = nodes
	workflow["links"] = [][]any{{1, 1, 1, 2, 0, "CLIP"}, {2, 1, 1, 3, 0, "CLIP"}, {3, 1, 0, 7, 0, "MODEL"}, {4, 7, 3, 8, 0, "MODEL"}, {5, 2, 0, 8, 1, "CONDITIONING"}, {6, 4, 0, 5, 0, "IMAGE"}, {7, 5, 0, 6, 0, "IMAGE"}, {8, 1, 2, 6, 1, "VAE"}, {9, 3, 0, 8, 2, "CONDITIONING"}, {10, 6, 0, 8, 3, "LATENT"}, {11, 1, 2, 8, 4, "VAE"}, {12, 7, 0, 2, 1, "STRING"}, {13, 7, 2, 3, 1, "STRING"}, {14, 7, 1, 8, 5, "STRING"}}
	return workflow
}

// fixtureQwenAnimationWorkflow retains native image-edit conditioning during sequence conversion.
func fixtureQwenAnimationWorkflow() map[string]any {
	workflow := fixtureAnimationWorkflow()
	var nodes []map[string]any
	for _, node := range workflow["nodes"].([]map[string]any) {
		id := node["id"].(int)
		if id == 2 || id == 3 {
			continue
		}
		if id == 8 {
			node["type"] = "ImagegenQwenChainedImageToImageGIF"
			node["inputs"] = []map[string]any{{"name": "model", "type": "MODEL", "link": 4}, {"name": "prompts", "type": "STRING", "link": 12}, {"name": "negative_prompt", "type": "STRING", "link": 13}, {"name": "latent_image", "type": "LATENT", "link": 10}, {"name": "vae", "type": "VAE", "link": 11}, {"name": "clip", "type": "CLIP", "link": 15}, {"name": "reference_image", "type": "IMAGE", "link": 16}, {"name": "filename_prefix", "type": "STRING", "link": 14, "widget": map[string]string{"name": "filename_prefix"}}}
		}
		nodes = append(nodes, node)
	}
	workflow["nodes"] = nodes
	workflow["links"] = [][]any{{3, 1, 0, 7, 0, "MODEL"}, {4, 7, 3, 8, 0, "MODEL"}, {6, 4, 0, 5, 0, "IMAGE"}, {7, 5, 0, 6, 0, "IMAGE"}, {8, 1, 2, 6, 1, "VAE"}, {10, 6, 0, 8, 3, "LATENT"}, {11, 1, 2, 8, 4, "VAE"}, {12, 7, 0, 8, 1, "STRING"}, {13, 7, 2, 8, 2, "STRING"}, {14, 7, 1, 8, 7, "STRING"}, {15, 1, 1, 8, 5, "CLIP"}, {16, 5, 0, 8, 6, "IMAGE"}}
	return workflow
}

// TestMain supplies a controlled upstream CLI process tree for cancellation acceptance.
func TestMain(m *testing.M) {
	mode := os.Getenv("EVALUATION_CLI_PROCESS_FIXTURE")
	if mode == "descendant" {
		time.Sleep(4 * time.Second)
		os.Exit(0)
	}
	if mode == "parent" {
		self, err := os.Executable()
		if err != nil {
			os.Exit(2)
		}
		child := exec.Command(self)
		child.Env = append(os.Environ(), "EVALUATION_CLI_PROCESS_FIXTURE=descendant")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		if err := os.WriteFile(os.Getenv("EVALUATION_CLI_PID_FILE"), []byte(strconv.Itoa(child.Process.Pid)), 0o600); err != nil {
			os.Exit(4)
		}
		if err := child.Wait(); err != nil {
			os.Exit(5)
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// TestCLICancellation terminates a stalled upstream CLI tree through the real bridge.
func TestCLICancellation(t *testing.T) {
	directory := t.TempDir()
	api := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `{}`) }))
	defer api.Close()
	ca := filepath.Join(directory, "ca.pem")
	if err := os.WriteFile(ca, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: api.Certificate().Raw}), 0o600); err != nil {
		t.Fatalf("write cancellation CA: %v", err)
	}
	bridge, err := runfiles.Rlocation(os.Getenv("BRIDGE_BINARY"))
	if err != nil {
		t.Fatalf("resolve cancellation bridge: %v", err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve controlled CLI: %v", err)
	}
	marker := filepath.Join(directory, "descendant.pid")
	t.Setenv("IMAGEGEN_WEB_USERNAME", "fixture")
	t.Setenv("IMAGEGEN_WEB_PASSWORD", "fixture-web-password")
	t.Setenv("EVALUATION_CLI_PROCESS_FIXTURE", "parent")
	t.Setenv("EVALUATION_CLI_PID_FILE", marker)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	started := time.Now()
	go func() {
		_, err := callCLI(ctx, Options{Bridge: bridge, CLI: self, CAFile: ca}, api.URL, "--json", "system-stats")
		done <- err
	}()
	deadline := time.Now().Add(3 * time.Second)
	var pid int
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(marker)
		if err == nil {
			pid, _ = strconv.Atoi(string(data))
			if pid > 0 {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	if pid <= 0 {
		cancel()
		<-done
		t.Fatal("controlled CLI did not start its stdout-holding descendant")
	}
	defer syscall.Kill(pid, syscall.SIGKILL)
	cancelled := time.Now()
	cancel()
	if err := <-done; err == nil {
		t.Fatal("cancelled bridge succeeded")
	}
	if time.Since(cancelled) > 1500*time.Millisecond {
		t.Fatal("CLI descendant retained the event pipe beyond the cancellation bound")
	}
	data, readErr := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if readErr == nil {
		_, after, _ := strings.Cut(string(data), ") ")
		if !strings.HasPrefix(after, "Z ") {
			t.Fatal("CLI descendant remained running after cancellation")
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		receipt := map[string]any{"process_tree_terminated": true, "elapsed_seconds": time.Since(started).Seconds(), "live_operations": 0}
		encoded, err := json.Marshal(receipt)
		if err != nil {
			t.Fatalf("encode cancellation receipt: %v", err)
		}
		if err := os.WriteFile(filepath.Join(output, "cli-cancellation.json"), encoded, 0o600); err != nil {
			t.Fatalf("save cancellation receipt: %v", err)
		}
	}
}
