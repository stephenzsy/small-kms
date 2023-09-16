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
  Certificate,
  CertificateEnrollRequest,
  CertificateRef,
} from '../models';
import {
    CertificateFromJSON,
    CertificateToJSON,
    CertificateEnrollRequestFromJSON,
    CertificateEnrollRequestToJSON,
    CertificateRefFromJSON,
    CertificateRefToJSON,
} from '../models';

export interface EnrollCertificateV1Request {
    certificateEnrollRequest: CertificateEnrollRequest;
}

export interface GetCertificateV1Request {
    namespaceId: string;
    id: string;
    accept?: GetCertificateV1AcceptEnum;
}

export interface ListCertificatesV1Request {
    namespaceId: string;
}

/**
 * 
 */
export class CertsApi extends runtime.BaseAPI {

    /**
     * Enroll certificate
     */
    async enrollCertificateV1Raw(requestParameters: EnrollCertificateV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Certificate>> {
        if (requestParameters.certificateEnrollRequest === null || requestParameters.certificateEnrollRequest === undefined) {
            throw new runtime.RequiredError('certificateEnrollRequest','Required parameter requestParameters.certificateEnrollRequest was null or undefined when calling enrollCertificateV1.');
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
            path: `/v1/certificates/enroll`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: CertificateEnrollRequestToJSON(requestParameters.certificateEnrollRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateFromJSON(jsonValue));
    }

    /**
     * Enroll certificate
     */
    async enrollCertificateV1(requestParameters: EnrollCertificateV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Certificate> {
        const response = await this.enrollCertificateV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get certificate
     */
    async getCertificateV1Raw(requestParameters: GetCertificateV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateRef>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getCertificateV1.');
        }

        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling getCertificateV1.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (requestParameters.accept !== undefined && requestParameters.accept !== null) {
            headerParameters['Accept'] = String(requestParameters.accept);
        }

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/{namespaceId}/certificates/{id}`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateRefFromJSON(jsonValue));
    }

    /**
     * Get certificate
     */
    async getCertificateV1(requestParameters: GetCertificateV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateRef> {
        const response = await this.getCertificateV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List certificates
     */
    async listCertificatesV1Raw(requestParameters: ListCertificatesV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<CertificateRef>>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling listCertificatesV1.');
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
            path: `/v1/{namespaceId}/certificates`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(CertificateRefFromJSON));
    }

    /**
     * List certificates
     */
    async listCertificatesV1(requestParameters: ListCertificatesV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<CertificateRef>> {
        const response = await this.listCertificatesV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

}

/**
 * @export
 */
export const GetCertificateV1AcceptEnum = {
    Accept_Json: 'application/json',
    Accept_Pem: 'application/x-pem-file',
    Accept_X509CaCert: 'application/x-x509-ca-cert'
} as const;
export type GetCertificateV1AcceptEnum = typeof GetCertificateV1AcceptEnum[keyof typeof GetCertificateV1AcceptEnum];
