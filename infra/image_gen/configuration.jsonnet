local models = import '../../projects/comfyui_tools/workflows/models.json';
local deployment = import 'deployment.json';
{ 'config.json': deployment + models }
