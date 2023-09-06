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
  NamespaceRef,
  NamespaceType,
  RegisterNamespaceV1200Response,
} from '../models';
import {
    NamespaceRefFromJSON,
    NamespaceRefToJSON,
    NamespaceTypeFromJSON,
    NamespaceTypeToJSON,
    RegisterNamespaceV1200ResponseFromJSON,
    RegisterNamespaceV1200ResponseToJSON,
} from '../models';

export interface ListNamespacesV1Request {
    namespaceType: NamespaceType;
}

export interface RegisterNamespaceV1Request {
    namespaceId: string;
}

/**
 * 
 */
export class DirectoryApi extends runtime.BaseAPI {

    /**
     * List namespaces
     */
    async listNamespacesV1Raw(requestParameters: ListNamespacesV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<NamespaceRef>>> {
        if (requestParameters.namespaceType === null || requestParameters.namespaceType === undefined) {
            throw new runtime.RequiredError('namespaceType','Required parameter requestParameters.namespaceType was null or undefined when calling listNamespacesV1.');
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
            path: `/v1/namespaces/{namespaceType}`.replace(`{${"namespaceType"}}`, encodeURIComponent(String(requestParameters.namespaceType))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(NamespaceRefFromJSON));
    }

    /**
     * List namespaces
     */
    async listNamespacesV1(requestParameters: ListNamespacesV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<NamespaceRef>> {
        const response = await this.listNamespacesV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Register namespace
     */
    async registerNamespaceV1Raw(requestParameters: RegisterNamespaceV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<RegisterNamespaceV1200Response>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling registerNamespaceV1.');
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
            path: `/v1/{namespaceId}/register`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => RegisterNamespaceV1200ResponseFromJSON(jsonValue));
    }

    /**
     * Register namespace
     */
    async registerNamespaceV1(requestParameters: RegisterNamespaceV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<RegisterNamespaceV1200Response> {
        const response = await this.registerNamespaceV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

}
