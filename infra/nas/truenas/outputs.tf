output "dataset_mountpoints" {
  description = "Managed filesystem dataset mountpoints."
  value       = { for name, dataset in truenas_dataset.managed : name => dataset.mountpoint }
}
