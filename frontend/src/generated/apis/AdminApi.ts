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
  CertificateInfo,
  CertificateTemplate,
  CertificateTemplateParameters,
  DeviceServicePrincipal,
  NamespaceTypeShortName,
  Ref,
} from '../models';
import {
    CertificateInfoFromJSON,
    CertificateInfoToJSON,
    CertificateTemplateFromJSON,
    CertificateTemplateToJSON,
    CertificateTemplateParametersFromJSON,
    CertificateTemplateParametersToJSON,
    DeviceServicePrincipalFromJSON,
    DeviceServicePrincipalToJSON,
    NamespaceTypeShortNameFromJSON,
    NamespaceTypeShortNameToJSON,
    RefFromJSON,
    RefToJSON,
} from '../models';

export interface GetCertificateTemplateV2Request {
    namespaceType: NamespaceTypeShortName;
    namespaceId: string;
    templateId: string;
}

export interface GetCertificateV2Request {
    namespaceType: NamespaceTypeShortName;
    namespaceId: string;
    templateId: string;
    certId: string;
    apply?: boolean;
    includeCertificate?: GetCertificateV2IncludeCertificateEnum;
}

export interface LinkDeviceServicePrincipalV2Request {
    namespaceId: string;
}

export interface ListCertificateTemplatesV2Request {
    namespaceType: NamespaceTypeShortName;
    namespaceId: string;
}

export interface ListNamespacesByTypeV2Request {
    namespaceType: NamespaceTypeShortName;
}

export interface PutCertificateTemplateV2Request {
    namespaceType: NamespaceTypeShortName;
    namespaceId: string;
    templateId: string;
    displayName: string;
    certificateTemplateParameters: CertificateTemplateParameters;
}

/**
 * 
 */
export class AdminApi extends runtime.BaseAPI {

    /**
     * Get certificate template
     */
    async getCertificateTemplateV2Raw(requestParameters: GetCertificateTemplateV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateTemplate>> {
        if (requestParameters.namespaceType === null || requestParameters.namespaceType === undefined) {
            throw new runtime.RequiredError('namespaceType','Required parameter requestParameters.namespaceType was null or undefined when calling getCertificateTemplateV2.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getCertificateTemplateV2.');
        }

        if (requestParameters.templateId === null || requestParameters.templateId === undefined) {
            throw new runtime.RequiredError('templateId','Required parameter requestParameters.templateId was null or undefined when calling getCertificateTemplateV2.');
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
            path: `/v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId}`.replace(`{${"namespaceType"}}`, encodeURIComponent(String(requestParameters.namespaceType))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"templateId"}}`, encodeURIComponent(String(requestParameters.templateId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateTemplateFromJSON(jsonValue));
    }

    /**
     * Get certificate template
     */
    async getCertificateTemplateV2(requestParameters: GetCertificateTemplateV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateTemplate> {
        const response = await this.getCertificateTemplateV2Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get certificate
     */
    async getCertificateV2Raw(requestParameters: GetCertificateV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateInfo>> {
        if (requestParameters.namespaceType === null || requestParameters.namespaceType === undefined) {
            throw new runtime.RequiredError('namespaceType','Required parameter requestParameters.namespaceType was null or undefined when calling getCertificateV2.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getCertificateV2.');
        }

        if (requestParameters.templateId === null || requestParameters.templateId === undefined) {
            throw new runtime.RequiredError('templateId','Required parameter requestParameters.templateId was null or undefined when calling getCertificateV2.');
        }

        if (requestParameters.certId === null || requestParameters.certId === undefined) {
            throw new runtime.RequiredError('certId','Required parameter requestParameters.certId was null or undefined when calling getCertificateV2.');
        }

        const queryParameters: any = {};

        if (requestParameters.apply !== undefined) {
            queryParameters['apply'] = requestParameters.apply;
        }

        if (requestParameters.includeCertificate !== undefined) {
            queryParameters['includeCertificate'] = requestParameters.includeCertificate;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId}/certificates/{certId}`.replace(`{${"namespaceType"}}`, encodeURIComponent(String(requestParameters.namespaceType))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"templateId"}}`, encodeURIComponent(String(requestParameters.templateId))).replace(`{${"certId"}}`, encodeURIComponent(String(requestParameters.certId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateInfoFromJSON(jsonValue));
    }

    /**
     * Get certificate
     */
    async getCertificateV2(requestParameters: GetCertificateV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateInfo> {
        const response = await this.getCertificateV2Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Link device service principal
     */
    async linkDeviceServicePrincipalV2Raw(requestParameters: LinkDeviceServicePrincipalV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<DeviceServicePrincipal>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling linkDeviceServicePrincipalV2.');
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
            path: `/v2/device/{namespaceId}/link-service-principal`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DeviceServicePrincipalFromJSON(jsonValue));
    }

    /**
     * Link device service principal
     */
    async linkDeviceServicePrincipalV2(requestParameters: LinkDeviceServicePrincipalV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<DeviceServicePrincipal> {
        const response = await this.linkDeviceServicePrincipalV2Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List certificate templates
     */
    async listCertificateTemplatesV2Raw(requestParameters: ListCertificateTemplatesV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<Ref>>> {
        if (requestParameters.namespaceType === null || requestParameters.namespaceType === undefined) {
            throw new runtime.RequiredError('namespaceType','Required parameter requestParameters.namespaceType was null or undefined when calling listCertificateTemplatesV2.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling listCertificateTemplatesV2.');
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
            path: `/v2/{namespaceType}/{namespaceId}/certificate-templates`.replace(`{${"namespaceType"}}`, encodeURIComponent(String(requestParameters.namespaceType))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RefFromJSON));
    }

    /**
     * List certificate templates
     */
    async listCertificateTemplatesV2(requestParameters: ListCertificateTemplatesV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<Ref>> {
        const response = await this.listCertificateTemplatesV2Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List namespaces by type
     */
    async listNamespacesByTypeV2Raw(requestParameters: ListNamespacesByTypeV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<Ref>>> {
        if (requestParameters.namespaceType === null || requestParameters.namespaceType === undefined) {
            throw new runtime.RequiredError('namespaceType','Required parameter requestParameters.namespaceType was null or undefined when calling listNamespacesByTypeV2.');
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
            path: `/v2/{namespaceType}`.replace(`{${"namespaceType"}}`, encodeURIComponent(String(requestParameters.namespaceType))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RefFromJSON));
    }

    /**
     * List namespaces by type
     */
    async listNamespacesByTypeV2(requestParameters: ListNamespacesByTypeV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<Ref>> {
        const response = await this.listNamespacesByTypeV2Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Put certificate template
     */
    async putCertificateTemplateV2Raw(requestParameters: PutCertificateTemplateV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificateTemplate>> {
        if (requestParameters.namespaceType === null || requestParameters.namespaceType === undefined) {
            throw new runtime.RequiredError('namespaceType','Required parameter requestParameters.namespaceType was null or undefined when calling putCertificateTemplateV2.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling putCertificateTemplateV2.');
        }

        if (requestParameters.templateId === null || requestParameters.templateId === undefined) {
            throw new runtime.RequiredError('templateId','Required parameter requestParameters.templateId was null or undefined when calling putCertificateTemplateV2.');
        }

        if (requestParameters.displayName === null || requestParameters.displayName === undefined) {
            throw new runtime.RequiredError('displayName','Required parameter requestParameters.displayName was null or undefined when calling putCertificateTemplateV2.');
        }

        if (requestParameters.certificateTemplateParameters === null || requestParameters.certificateTemplateParameters === undefined) {
            throw new runtime.RequiredError('certificateTemplateParameters','Required parameter requestParameters.certificateTemplateParameters was null or undefined when calling putCertificateTemplateV2.');
        }

        const queryParameters: any = {};

        if (requestParameters.displayName !== undefined) {
            queryParameters['displayName'] = requestParameters.displayName;
        }

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
            path: `/v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId}`.replace(`{${"namespaceType"}}`, encodeURIComponent(String(requestParameters.namespaceType))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"templateId"}}`, encodeURIComponent(String(requestParameters.templateId))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: CertificateTemplateParametersToJSON(requestParameters.certificateTemplateParameters),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificateTemplateFromJSON(jsonValue));
    }

    /**
     * Put certificate template
     */
    async putCertificateTemplateV2(requestParameters: PutCertificateTemplateV2Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificateTemplate> {
        const response = await this.putCertificateTemplateV2Raw(requestParameters, initOverrides);
        return await response.value();
    }

}

/**
 * @export
 */
export const GetCertificateV2IncludeCertificateEnum = {
    IncludeJWK: 'jwk',
    IncludePEM: 'pem'
} as const;
export type GetCertificateV2IncludeCertificateEnum = typeof GetCertificateV2IncludeCertificateEnum[keyof typeof GetCertificateV2IncludeCertificateEnum];
