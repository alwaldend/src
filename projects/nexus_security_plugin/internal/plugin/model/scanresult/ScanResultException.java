package com.alwaldend.src.projects.nexus_security_plugin.model.scanresult;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.io.Serializable;
import org.apache.commons.lang.exception.ExceptionUtils;

public class ScanResultException implements Serializable {

  @JsonProperty("class")
  private final Class<? extends Throwable> exceptionClass;

  @JsonProperty("stack_trace")
  private final String stackTrace;

  public ScanResultException(Throwable exception) {
    this.exceptionClass = exception.getClass();
    this.stackTrace = ExceptionUtils.getStackTrace(exception);
  }
}
