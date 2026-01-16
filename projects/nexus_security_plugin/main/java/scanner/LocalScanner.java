package com.alwaldend.src.projects.nexus_security_plugin.scanner;

import com.alwaldend.src.projects.nexus_security_plugin.main.BundleHelper;
import com.alwaldend.src.projects.nexus_security_plugin.model.bundle.configuration.BundleConfiguration;
import com.alwaldend.src.projects.nexus_security_plugin.model.bundle.configuration.BundleConfigurationScannerLocal;
import com.alwaldend.src.projects.nexus_security_plugin.model.information.request.RequestInformation;
import com.alwaldend.src.projects.nexus_security_plugin.model.scanresult.ScanResult;
import com.alwaldend.src.projects.nexus_security_plugin.model.scanresult.ScanResultType;
import com.alwaldend.src.projects.nexus_security_plugin.model.white_list.WhiteList;
import org.joda.time.DateTime;
import org.sonatype.nexus.repository.storage.AssetStore;

import javax.inject.Inject;
import javax.inject.Named;

@Named
public class LocalScanner extends AbstractScanner {

    @Inject
    public LocalScanner(AssetStore assetStore,
                        BundleHelper bundleHelper) {
        super(assetStore, bundleHelper);
    }

    @Override
    ScanResult scanImpl(RequestInformation information) {
        BundleConfiguration config = bundleHelper
            .getBundleConfiguration();
        BundleConfigurationScannerLocal scannerConfig = config
            .getScanners()
            .getLocal();
        ScanResultType scanResultType = scannerConfig
            .getWhiteList()
            .contains(information);
        if (!WhiteList.isFailure(scanResultType)) {
            return new ScanResult(true, scanResultType);
        }
        DateTime lastModified = information.getComponent().getLastModified();
        assert lastModified != null;
        if (lastModified.isBefore(scannerConfig.getLastModified())) {
            return new ScanResult(true, ScanResultType.LAST_MODIFIED_VALID);
        }
        return new ScanResult(false, ScanResultType.LAST_MODIFIED_INVALID);
    }

    @Override
    public boolean failOnErrors() {
        return bundleHelper
            .getBundleConfiguration()
            .getScanners()
            .getLocal()
            .isFailOnErrors();
    }
}
