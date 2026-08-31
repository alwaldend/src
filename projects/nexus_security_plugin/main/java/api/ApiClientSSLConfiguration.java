package com.alwaldend.src.projects.nexus_security_plugin.api;

import java.security.cert.X509Certificate;
import javax.annotation.Nonnull;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

public final class ApiClientSSLConfiguration {

  private ApiClientSSLConfiguration() {}

  @Nonnull
  public static TrustManager[] buildUnsafeTrustManager() {
    return new TrustManager[] {
      new X509TrustManager() {
        @SuppressWarnings("squid:S4424")
        @Override
        public void checkClientTrusted(X509Certificate[] x509Certificates, String authType) {}

        @SuppressWarnings("squid:S4424")
        @Override
        public void checkServerTrusted(X509Certificate[] x509Certificates, String authType) {}

        @Override
        public X509Certificate[] getAcceptedIssuers() {
          return new X509Certificate[0];
        }
      }
    };
  }
}
