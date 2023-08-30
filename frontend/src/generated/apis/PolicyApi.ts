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
  CertificateEnrollmentParameters,
  CertificateEnrollmentPolicy,
} from '../models';
import {
    CertificateEnrollmentParametersFromJSON,
    CertificateEnrollmentParametersToJSON,
    CertificateEnrollmentPolicyFromJSON,
    CertificateEnrollmentPolicyToJSON,
} from '../models';

export interface GetPolicyCertEnrollV1Request {
    namespaceId: string;
}

export interface PutPolicyCertEnrollV1Request {
    namespaceId: string;
    certificateEnrollmentParameters: CertificateEnrollmentParameters;
}

/**
 * 
 */
export class PolicyApi extends runtime.BaseAPI {

    /**
     * Get Certificate Enrollment Policy
     */
    async getPolicyCertEnrollV1Raw(requestParameters: GetPolicyCertEnrollV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateEnrollmentPolicy>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getPolicyCertEnrollV1.');
        }

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
            path: `/v1/{namespaceId}/policy/cert-enrollment`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateEnrollmentPolicyFromJSON(jsonValue));
    }

    /**
     * Get Certificate Enrollment Policy
     */
    async getPolicyCertEnrollV1(requestParameters: GetPolicyCertEnrollV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateEnrollmentPolicy> {
        const response = await this.getPolicyCertEnrollV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Put Certificate Enrollment Policy
     */
    async putPolicyCertEnrollV1Raw(requestParameters: PutPolicyCertEnrollV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateEnrollmentPolicy>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling putPolicyCertEnrollV1.');
        }

        if (requestParameters.certificateEnrollmentParameters === null || requestParameters.certificateEnrollmentParameters === undefined) {
            throw new runtime.RequiredError('certificateEnrollmentParameters','Required parameter requestParameters.certificateEnrollmentParameters was null or undefined when calling putPolicyCertEnrollV1.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/{namespaceId}/policy/cert-enrollment`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: CertificateEnrollmentParametersToJSON(requestParameters.certificateEnrollmentParameters),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateEnrollmentPolicyFromJSON(jsonValue));
    }

    /**
     * Put Certificate Enrollment Policy
     */
    async putPolicyCertEnrollV1(requestParameters: PutPolicyCertEnrollV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateEnrollmentPolicy> {
        const response = await this.putPolicyCertEnrollV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

}
