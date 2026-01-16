package com.alwaldend.src.projects.nexus_security_plugin.api;

import com.alwaldend.src.projects.nexus_security_plugin.model.information.request.RequestInformation;
import com.alwaldend.src.projects.nexus_security_plugin.model.scanresult.ScanResult;
import retrofit2.Call;
import retrofit2.http.Body;
import retrofit2.http.PUT;

public interface RemoteScanApi {

    @PUT("api/v1/check")
    Call<ScanResult> check(@Body RequestInformation body);

}
