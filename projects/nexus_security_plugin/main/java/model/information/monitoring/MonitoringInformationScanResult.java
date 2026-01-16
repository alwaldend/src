package com.alwaldend.src.projects.nexus_security_plugin.model.information.monitoring;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.alwaldend.src.projects.nexus_security_plugin.model.scanresult.ScanResult;
import com.alwaldend.src.projects.nexus_security_plugin.scanner.AbstractScanner;

import java.io.Serializable;

public class MonitoringInformationScanResult implements Serializable {

    @JsonProperty
    private final Class<? extends AbstractScanner> type;

    @JsonProperty
    private final ScanResult result;

    public MonitoringInformationScanResult(Class<? extends AbstractScanner> type,
                                           ScanResult result) {
        this.type = type;
        this.result = result;
    }

    public ScanResult getResult() {
        return result;
    }
}
