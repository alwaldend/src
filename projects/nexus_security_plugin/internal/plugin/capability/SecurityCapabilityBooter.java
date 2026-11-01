package com.alwaldend.src.projects.nexus_security_plugin.capability;

import static org.sonatype.nexus.common.app.ManagedLifecycle.Phase.CAPABILITIES;

import javax.inject.Named;
import javax.inject.Singleton;
import org.sonatype.nexus.capability.CapabilityBooterSupport;
import org.sonatype.nexus.capability.CapabilityRegistry;
import org.sonatype.nexus.common.app.ManagedLifecycle;

@Named
@Singleton
@ManagedLifecycle(phase = CAPABILITIES)
public class SecurityCapabilityBooter extends CapabilityBooterSupport {
  @Override
  protected void boot(final CapabilityRegistry registry) throws Exception {
    maybeAddCapability(
        registry,
        SecurityCapabilityDescriptor.CAPABILITY_TYPE,
        true,
        null,
        SecurityCapabilityFormField.createDefaultProperties());
  }
}
