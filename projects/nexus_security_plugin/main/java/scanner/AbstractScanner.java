package com.alwaldend.src.projects.nexus_security_plugin.scanner;

import com.alwaldend.src.projects.nexus_security_plugin.logging.SecurityLogConfiguration;
import com.alwaldend.src.projects.nexus_security_plugin.main.BundleHelper;
import com.alwaldend.src.projects.nexus_security_plugin.model.information.request.RequestInformation;
import com.alwaldend.src.projects.nexus_security_plugin.model.scanresult.ScanResult;
import com.alwaldend.src.projects.nexus_security_plugin.model.scanresult.ScanResultType;
import org.sonatype.nexus.repository.storage.AssetStore;

abstract public class AbstractScanner {

    final AssetStore assetStore;
    final BundleHelper bundleHelper;

    public AbstractScanner(final AssetStore assetStore,
                           final BundleHelper bundleHelper) {
        this.assetStore = assetStore;
        this.bundleHelper = bundleHelper;
    }

    public ScanResult scan(RequestInformation information) {
        try {
            return scanImpl(information);
        } catch (Throwable exception) {
            boolean fail = failOnErrors();
            SecurityLogConfiguration.LOG.error("Could not scan asset {}",
                information.getRequest().getPath(),
                exception);
            return new ScanResult(
                !fail,
                fail ? ScanResultType.EXCEPTION
                    : ScanResultType.EXCEPTION_IGNORED)
                .setException(exception);
        }
    }

    abstract ScanResult scanImpl(RequestInformation information) throws Throwable;

    abstract boolean failOnErrors();
}
