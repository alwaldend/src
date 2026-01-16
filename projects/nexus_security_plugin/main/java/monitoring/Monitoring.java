package com.alwaldend.src.projects.nexus_security_plugin.monitoring;

import com.alwaldend.src.projects.nexus_security_plugin.api.MonitoringApi;
import com.alwaldend.src.projects.nexus_security_plugin.logging.SecurityLogConfiguration;
import com.alwaldend.src.projects.nexus_security_plugin.main.BundleHelper;
import com.alwaldend.src.projects.nexus_security_plugin.model.bundle.configuration.BundleConfiguration;
import com.alwaldend.src.projects.nexus_security_plugin.model.information.monitoring.MonitoringInformation;

import javax.inject.Inject;
import javax.inject.Named;
import java.io.IOException;
import java.io.Serializable;

@Named
public class Monitoring implements Serializable {


    private final BundleHelper bundleHelper;

    @Inject
    public Monitoring(final BundleHelper bundleHelper) {
        this.bundleHelper = bundleHelper;
    }

    public void send(MonitoringInformation information) {
        BundleConfiguration config = bundleHelper
            .getBundleConfiguration();
        if (!config.getMonitoring().isEnabled()) {
            return;
        }
        try {
            sendImpl(information, config);
        } catch (Throwable exception) {
            SecurityLogConfiguration.LOG.error("Could not send monitoring information", exception);
        }
    }

    public void sendImpl(MonitoringInformation information,
                         BundleConfiguration config) throws IOException {
        MonitoringApi api = bundleHelper.getMonitoringApi();

        api.bulk("{\"index\":{ } }\n"
                + BundleHelper.MAPPER_JSON
                .writeValueAsString(information),
            config.getMonitoring().getIndex(),
            BundleHelper.todayDate(),
            config.getMonitoring().getPipeline()
        ).execute();
    }
}
