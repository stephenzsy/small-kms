/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  RequestDiagnostics,
} from '../models';
import {
    RequestDiagnosticsFromJSON,
    RequestDiagnosticsToJSON,
} from '../models';

/**
 * 
 */
export class DiagnosticsApi extends runtime.BaseAPI {

    /**
     * Get diagnostics
     */
    async getDiagnosticsV1Raw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<RequestDiagnostics>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/diagnostics`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => RequestDiagnosticsFromJSON(jsonValue));
    }

    /**
     * Get diagnostics
     */
    async getDiagnosticsV1(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<RequestDiagnostics> {
        const response = await this.getDiagnosticsV1Raw(initOverrides);
        return await response.value();
    }

}
