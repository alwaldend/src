// Model filenames come from their existing public configuration owner.
local configuration = import '../../../infra/image_gen/config.json';
local g = import 'internal/graph.libsonnet';
// Recorded workflow defaults remain usable when native downloading is disabled.
local defaults = {
  qwen: {
    checkpoint: 'qwen_image_edit_2511_fp8mixed.safetensors',
    text_encoder: 'qwen_2.5_vl_7b_fp8_scaled.safetensors',
    vae: 'qwen_image_vae.safetensors',
  },
  flux: {
    checkpoint: 'flux-2-klein-base-9b-fp8.safetensors',
    text_encoder: 'qwen_3_8b_fp8mixed.safetensors',
    vae: 'full_encoder_small_decoder.safetensors',
  },
  juggernaut: {
    checkpoint: 'Juggernaut-XL_v9_RunDiffusionPhoto_v2.safetensors',
    ipadapter: 'ip-adapter-plus_sdxl_vit-h.safetensors',
    clip_vision: 'CLIP-ViT-H-14-laion2B-s32B-b79K.safetensors',
  },
  'flux2-dev': {
    checkpoint: 'flux2_dev_fp8mixed.safetensors',
    text_encoder: 'mistral_3_small_flux2_fp8.safetensors',
    vae: 'flux2-vae.safetensors',
  },
};
local downloadKind(entry) =
  if std.get(entry, 'kind', null) == null then 'checkpoint' else entry.kind;
// First matching files select workflow components; additional downloads remain available.
local filenames(profile) = {
  [kind]:
    local matches = [
      entry
      for entry in configuration.model_profiles[profile].downloads
      if downloadKind(entry) == kind
    ];
    if std.length(matches) == 0 then defaults[profile][kind] else matches[0].filename
  for kind in std.objectFields(defaults[profile])
};
local families = {
  flux: (import 'flux/workflows.libsonnet')(filenames('flux')),
  qwen: (import 'qwen/workflows.libsonnet')(filenames('qwen')),
  juggernaut: (import 'juggernaut/workflows.libsonnet')(filenames('juggernaut')),
  chroma: (import 'chroma/workflows.libsonnet')({}),
  'flux2-dev': (import 'flux2-dev/workflows.libsonnet')(filenames('flux2-dev')),
};
{
  [profile + '/' + profile + '-' + name + '.json']: g.render(families[profile][name])
  for profile in std.objectFields(families)
  for name in std.objectFields(families[profile])
}
