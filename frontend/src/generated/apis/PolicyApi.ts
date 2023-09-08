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
  ApplyPolicyRequest,
  GetPolicyV1404Response,
  Policy,
  PolicyParameters,
  PolicyRef,
  PolicyState,
} from '../models';
import {
    ApplyPolicyRequestFromJSON,
    ApplyPolicyRequestToJSON,
    GetPolicyV1404ResponseFromJSON,
    GetPolicyV1404ResponseToJSON,
    PolicyFromJSON,
    PolicyToJSON,
    PolicyParametersFromJSON,
    PolicyParametersToJSON,
    PolicyRefFromJSON,
    PolicyRefToJSON,
    PolicyStateFromJSON,
    PolicyStateToJSON,
} from '../models';

export interface ApplyPolicyV1Request {
    namespaceId: string;
    policyId: string;
    applyPolicyRequest?: ApplyPolicyRequest;
}

export interface GetPolicyV1Request {
    namespaceId: string;
    policyId: string;
}

export interface ListPoliciesV1Request {
    namespaceId: string;
}

export interface PutPolicyV1Request {
    namespaceId: string;
    policyId: string;
    policyParameters: PolicyParameters;
}

/**
 * 
 */
export class PolicyApi extends runtime.BaseAPI {

    /**
     * Apply policy
     */
    async applyPolicyV1Raw(requestParameters: ApplyPolicyV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<PolicyState>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling applyPolicyV1.');
        }

        if (requestParameters.policyId === null || requestParameters.policyId === undefined) {
            throw new runtime.RequiredError('policyId','Required parameter requestParameters.policyId was null or undefined when calling applyPolicyV1.');
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
            path: `/v1/{namespaceId}/policies/{policyId}/apply`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"policyId"}}`, encodeURIComponent(String(requestParameters.policyId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: ApplyPolicyRequestToJSON(requestParameters.applyPolicyRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => PolicyStateFromJSON(jsonValue));
    }

    /**
     * Apply policy
     */
    async applyPolicyV1(requestParameters: ApplyPolicyV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<PolicyState> {
        const response = await this.applyPolicyV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get Certificate Policy
     */
    async getPolicyV1Raw(requestParameters: GetPolicyV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Policy>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getPolicyV1.');
        }

        if (requestParameters.policyId === null || requestParameters.policyId === undefined) {
            throw new runtime.RequiredError('policyId','Required parameter requestParameters.policyId was null or undefined when calling getPolicyV1.');
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
            path: `/v1/{namespaceId}/policies/{policyId}`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"policyId"}}`, encodeURIComponent(String(requestParameters.policyId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => PolicyFromJSON(jsonValue));
    }

    /**
     * Get Certificate Policy
     */
    async getPolicyV1(requestParameters: GetPolicyV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Policy> {
        const response = await this.getPolicyV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List policies
     */
    async listPoliciesV1Raw(requestParameters: ListPoliciesV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<PolicyRef>>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling listPoliciesV1.');
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
            path: `/v1/{namespaceId}/policies`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(PolicyRefFromJSON));
    }

    /**
     * List policies
     */
    async listPoliciesV1(requestParameters: ListPoliciesV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<PolicyRef>> {
        const response = await this.listPoliciesV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Put Policy
     */
    async putPolicyV1Raw(requestParameters: PutPolicyV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Policy>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling putPolicyV1.');
        }

        if (requestParameters.policyId === null || requestParameters.policyId === undefined) {
            throw new runtime.RequiredError('policyId','Required parameter requestParameters.policyId was null or undefined when calling putPolicyV1.');
        }

        if (requestParameters.policyParameters === null || requestParameters.policyParameters === undefined) {
            throw new runtime.RequiredError('policyParameters','Required parameter requestParameters.policyParameters was null or undefined when calling putPolicyV1.');
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
            path: `/v1/{namespaceId}/policies/{policyId}`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"policyId"}}`, encodeURIComponent(String(requestParameters.policyId))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: PolicyParametersToJSON(requestParameters.policyParameters),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => PolicyFromJSON(jsonValue));
    }

    /**
     * Put Policy
     */
    async putPolicyV1(requestParameters: PutPolicyV1Request, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Policy> {
        const response = await this.putPolicyV1Raw(requestParameters, initOverrides);
        return await response.value();
    }

}