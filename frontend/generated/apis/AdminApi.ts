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
  CertificateCategory,
  CertificateFileFormat,
  CertificateRef,
  CreateCertificateParameters,
} from '../models';
import {
    CertificateCategoryFromJSON,
    CertificateCategoryToJSON,
    CertificateFileFormatFromJSON,
    CertificateFileFormatToJSON,
    CertificateRefFromJSON,
    CertificateRefToJSON,
    CreateCertificateParametersFromJSON,
    CreateCertificateParametersToJSON,
} from '../models';

export interface CreateCertificateRequest {
    xMsClientPrincipalName: string;
    xMsClientPrincipalId: string;
    xMsClientRoles: string;
    force?: boolean;
    createCertificateParameters?: CreateCertificateParameters;
}

export interface DownloadCertificateRequest {
    xMsClientPrincipalName: string;
    xMsClientPrincipalId: string;
    xMsClientRoles: string;
    id: string;
    format?: CertificateFileFormat;
}

export interface ListCertificatesRequest {
    xMsClientPrincipalName: string;
    xMsClientPrincipalId: string;
    xMsClientRoles: string;
    category: CertificateCategory;
}

/**
 * 
 */
export class AdminApi extends runtime.BaseAPI {

    /**
     * Create certificate
     */
    async createCertificateRaw(requestParameters: CreateCertificateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateRef>> {
        if (requestParameters.xMsClientPrincipalName === null || requestParameters.xMsClientPrincipalName === undefined) {
            throw new runtime.RequiredError('xMsClientPrincipalName','Required parameter requestParameters.xMsClientPrincipalName was null or undefined when calling createCertificate.');
        }

        if (requestParameters.xMsClientPrincipalId === null || requestParameters.xMsClientPrincipalId === undefined) {
            throw new runtime.RequiredError('xMsClientPrincipalId','Required parameter requestParameters.xMsClientPrincipalId was null or undefined when calling createCertificate.');
        }

        if (requestParameters.xMsClientRoles === null || requestParameters.xMsClientRoles === undefined) {
            throw new runtime.RequiredError('xMsClientRoles','Required parameter requestParameters.xMsClientRoles was null or undefined when calling createCertificate.');
        }

        const queryParameters: any = {};

        if (requestParameters.force !== undefined) {
            queryParameters['force'] = requestParameters.force;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (requestParameters.xMsClientPrincipalName !== undefined && requestParameters.xMsClientPrincipalName !== null) {
            headerParameters['X-Ms-Client-Principal-Name'] = String(requestParameters.xMsClientPrincipalName);
        }

        if (requestParameters.xMsClientPrincipalId !== undefined && requestParameters.xMsClientPrincipalId !== null) {
            headerParameters['X-Ms-Client-Principal-Id'] = String(requestParameters.xMsClientPrincipalId);
        }

        if (requestParameters.xMsClientRoles !== undefined && requestParameters.xMsClientRoles !== null) {
            headerParameters['X-Ms-Client-Roles'] = String(requestParameters.xMsClientRoles);
        }

        const response = await this.request({
            path: `/admin/certificate`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: CreateCertificateParametersToJSON(requestParameters.createCertificateParameters),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateRefFromJSON(jsonValue));
    }

    /**
     * Create certificate
     */
    async createCertificate(requestParameters: CreateCertificateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateRef> {
        const response = await this.createCertificateRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Download certificate
     */
    async downloadCertificateRaw(requestParameters: DownloadCertificateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Blob>> {
        if (requestParameters.xMsClientPrincipalName === null || requestParameters.xMsClientPrincipalName === undefined) {
            throw new runtime.RequiredError('xMsClientPrincipalName','Required parameter requestParameters.xMsClientPrincipalName was null or undefined when calling downloadCertificate.');
        }

        if (requestParameters.xMsClientPrincipalId === null || requestParameters.xMsClientPrincipalId === undefined) {
            throw new runtime.RequiredError('xMsClientPrincipalId','Required parameter requestParameters.xMsClientPrincipalId was null or undefined when calling downloadCertificate.');
        }

        if (requestParameters.xMsClientRoles === null || requestParameters.xMsClientRoles === undefined) {
            throw new runtime.RequiredError('xMsClientRoles','Required parameter requestParameters.xMsClientRoles was null or undefined when calling downloadCertificate.');
        }

        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling downloadCertificate.');
        }

        const queryParameters: any = {};

        if (requestParameters.format !== undefined) {
            queryParameters['format'] = requestParameters.format;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (requestParameters.xMsClientPrincipalName !== undefined && requestParameters.xMsClientPrincipalName !== null) {
            headerParameters['X-Ms-Client-Principal-Name'] = String(requestParameters.xMsClientPrincipalName);
        }

        if (requestParameters.xMsClientPrincipalId !== undefined && requestParameters.xMsClientPrincipalId !== null) {
            headerParameters['X-Ms-Client-Principal-Id'] = String(requestParameters.xMsClientPrincipalId);
        }

        if (requestParameters.xMsClientRoles !== undefined && requestParameters.xMsClientRoles !== null) {
            headerParameters['X-Ms-Client-Roles'] = String(requestParameters.xMsClientRoles);
        }

        const response = await this.request({
            path: `/admin/certificate/{id}/download`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.BlobApiResponse(response);
    }

    /**
     * Download certificate
     */
    async downloadCertificate(requestParameters: DownloadCertificateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Blob> {
        const response = await this.downloadCertificateRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List certificates
     */
    async listCertificatesRaw(requestParameters: ListCertificatesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<CertificateRef>>> {
        if (requestParameters.xMsClientPrincipalName === null || requestParameters.xMsClientPrincipalName === undefined) {
            throw new runtime.RequiredError('xMsClientPrincipalName','Required parameter requestParameters.xMsClientPrincipalName was null or undefined when calling listCertificates.');
        }

        if (requestParameters.xMsClientPrincipalId === null || requestParameters.xMsClientPrincipalId === undefined) {
            throw new runtime.RequiredError('xMsClientPrincipalId','Required parameter requestParameters.xMsClientPrincipalId was null or undefined when calling listCertificates.');
        }

        if (requestParameters.xMsClientRoles === null || requestParameters.xMsClientRoles === undefined) {
            throw new runtime.RequiredError('xMsClientRoles','Required parameter requestParameters.xMsClientRoles was null or undefined when calling listCertificates.');
        }

        if (requestParameters.category === null || requestParameters.category === undefined) {
            throw new runtime.RequiredError('category','Required parameter requestParameters.category was null or undefined when calling listCertificates.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (requestParameters.xMsClientPrincipalName !== undefined && requestParameters.xMsClientPrincipalName !== null) {
            headerParameters['X-Ms-Client-Principal-Name'] = String(requestParameters.xMsClientPrincipalName);
        }

        if (requestParameters.xMsClientPrincipalId !== undefined && requestParameters.xMsClientPrincipalId !== null) {
            headerParameters['X-Ms-Client-Principal-Id'] = String(requestParameters.xMsClientPrincipalId);
        }

        if (requestParameters.xMsClientRoles !== undefined && requestParameters.xMsClientRoles !== null) {
            headerParameters['X-Ms-Client-Roles'] = String(requestParameters.xMsClientRoles);
        }

        const response = await this.request({
            path: `/admin/certificates/{category}`.replace(`{${"category"}}`, encodeURIComponent(String(requestParameters.category))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(CertificateRefFromJSON));
    }

    /**
     * List certificates
     */
    async listCertificates(requestParameters: ListCertificatesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<CertificateRef>> {
        const response = await this.listCertificatesRaw(requestParameters, initOverrides);
        return await response.value();
    }

}