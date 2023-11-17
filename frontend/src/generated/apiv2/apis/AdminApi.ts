/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  AgentConfig,
  AgentConfigFields,
  ApplicationByAppId,
  CertificatePolicy,
  CreateAgentRequest,
  CreateCertificatePolicyRequest,
  ErrorResult,
  KeyPolicy,
  NamespaceProvider,
  Ref,
} from '../models/index';
import {
    AgentConfigFromJSON,
    AgentConfigToJSON,
    AgentConfigFieldsFromJSON,
    AgentConfigFieldsToJSON,
    ApplicationByAppIdFromJSON,
    ApplicationByAppIdToJSON,
    CertificatePolicyFromJSON,
    CertificatePolicyToJSON,
    CreateAgentRequestFromJSON,
    CreateAgentRequestToJSON,
    CreateCertificatePolicyRequestFromJSON,
    CreateCertificatePolicyRequestToJSON,
    ErrorResultFromJSON,
    ErrorResultToJSON,
    KeyPolicyFromJSON,
    KeyPolicyToJSON,
    NamespaceProviderFromJSON,
    NamespaceProviderToJSON,
    RefFromJSON,
    RefToJSON,
} from '../models/index';

export interface CreateAgentOperationRequest {
    createAgentRequest?: CreateAgentRequest;
}

export interface GetAgentRequest {
    id: string;
}

export interface GetAgentConfigRequest {
    namespaceId: string;
}

export interface GetCertificatePolicyRequest {
    namespaceProvider: NamespaceProvider;
    namespaceId: string;
    id: string;
}

export interface GetKeyPolicyRequest {
    namespaceProvider: NamespaceProvider;
    namespaceId: string;
    id: string;
}

export interface GetSystemAppRequest {
    id: string;
}

export interface ListCertificatePoliciesRequest {
    namespaceProvider: NamespaceProvider;
    namespaceId: string;
}

export interface ListKeyPoliciesRequest {
    namespaceProvider: NamespaceProvider;
    namespaceId: string;
}

export interface ListProfilesRequest {
    namespaceProvider: NamespaceProvider;
}

export interface PutAgentConfigRequest {
    namespaceId: string;
    agentConfigFields: AgentConfigFields;
}

export interface PutCertificatePolicyRequest {
    namespaceProvider: NamespaceProvider;
    namespaceId: string;
    id: string;
    createCertificatePolicyRequest: CreateCertificatePolicyRequest;
}

export interface PutKeyPolicyRequest {
    namespaceProvider: NamespaceProvider;
    namespaceId: string;
    id: string;
}

export interface SyncSystemAppRequest {
    id: string;
}

/**
 * 
 */
export class AdminApi extends runtime.BaseAPI {

    /**
     * Create agent
     */
    async createAgentRaw(requestParameters: CreateAgentOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ApplicationByAppId>> {
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
            path: `/v2/agents`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: CreateAgentRequestToJSON(requestParameters.createAgentRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ApplicationByAppIdFromJSON(jsonValue));
    }

    /**
     * Create agent
     */
    async createAgent(requestParameters: CreateAgentOperationRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ApplicationByAppId> {
        const response = await this.createAgentRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get agent
     */
    async getAgentRaw(requestParameters: GetAgentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ApplicationByAppId>> {
        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling getAgent.');
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
            path: `/v2/agents/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ApplicationByAppIdFromJSON(jsonValue));
    }

    /**
     * Get agent
     */
    async getAgent(requestParameters: GetAgentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ApplicationByAppId> {
        const response = await this.getAgentRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get agent config
     */
    async getAgentConfigRaw(requestParameters: GetAgentConfigRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AgentConfig>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getAgentConfig.');
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
            path: `/v2/service-principal/{namespaceId}/agent-config`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AgentConfigFromJSON(jsonValue));
    }

    /**
     * Get agent config
     */
    async getAgentConfig(requestParameters: GetAgentConfigRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AgentConfig> {
        const response = await this.getAgentConfigRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get certificate policy
     */
    async getCertificatePolicyRaw(requestParameters: GetCertificatePolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificatePolicy>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling getCertificatePolicy.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getCertificatePolicy.');
        }

        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling getCertificatePolicy.');
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
            path: `/v2/{namespaceProvider}/{namespaceId}/certificate-policies/{id}`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificatePolicyFromJSON(jsonValue));
    }

    /**
     * Get certificate policy
     */
    async getCertificatePolicy(requestParameters: GetCertificatePolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificatePolicy> {
        const response = await this.getCertificatePolicyRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get key policy
     */
    async getKeyPolicyRaw(requestParameters: GetKeyPolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<KeyPolicy>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling getKeyPolicy.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling getKeyPolicy.');
        }

        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling getKeyPolicy.');
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
            path: `/v2/{namespaceProvider}/{namespaceId}/key-policies/{id}`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => KeyPolicyFromJSON(jsonValue));
    }

    /**
     * Get key policy
     */
    async getKeyPolicy(requestParameters: GetKeyPolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<KeyPolicy> {
        const response = await this.getKeyPolicyRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get system app
     */
    async getSystemAppRaw(requestParameters: GetSystemAppRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ApplicationByAppId>> {
        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling getSystemApp.');
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
            path: `/v2/system-apps/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ApplicationByAppIdFromJSON(jsonValue));
    }

    /**
     * Get system app
     */
    async getSystemApp(requestParameters: GetSystemAppRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ApplicationByAppId> {
        const response = await this.getSystemAppRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List certificate policies
     */
    async listCertificatePoliciesRaw(requestParameters: ListCertificatePoliciesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<Ref>>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling listCertificatePolicies.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling listCertificatePolicies.');
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
            path: `/v2/{namespaceProvider}/{namespaceId}/certificate-policies`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RefFromJSON));
    }

    /**
     * List certificate policies
     */
    async listCertificatePolicies(requestParameters: ListCertificatePoliciesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<Ref>> {
        const response = await this.listCertificatePoliciesRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * List key policies
     */
    async listKeyPoliciesRaw(requestParameters: ListKeyPoliciesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<Ref>>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling listKeyPolicies.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling listKeyPolicies.');
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
            path: `/v2/{namespaceProvider}/{namespaceId}/key-policies`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RefFromJSON));
    }

    /**
     * List key policies
     */
    async listKeyPolicies(requestParameters: ListKeyPoliciesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<Ref>> {
        const response = await this.listKeyPoliciesRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * list profiles
     */
    async listProfilesRaw(requestParameters: ListProfilesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<Ref>>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling listProfiles.');
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
            path: `/v2/profiles/{namespaceProvider}`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RefFromJSON));
    }

    /**
     * list profiles
     */
    async listProfiles(requestParameters: ListProfilesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<Ref>> {
        const response = await this.listProfilesRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Put agent config
     */
    async putAgentConfigRaw(requestParameters: PutAgentConfigRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AgentConfig>> {
        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling putAgentConfig.');
        }

        if (requestParameters.agentConfigFields === null || requestParameters.agentConfigFields === undefined) {
            throw new runtime.RequiredError('agentConfigFields','Required parameter requestParameters.agentConfigFields was null or undefined when calling putAgentConfig.');
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
            path: `/v2/service-principal/{namespaceId}/agent-config`.replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: AgentConfigFieldsToJSON(requestParameters.agentConfigFields),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AgentConfigFromJSON(jsonValue));
    }

    /**
     * Put agent config
     */
    async putAgentConfig(requestParameters: PutAgentConfigRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AgentConfig> {
        const response = await this.putAgentConfigRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * put certificate policy
     */
    async putCertificatePolicyRaw(requestParameters: PutCertificatePolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CertificatePolicy>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling putCertificatePolicy.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling putCertificatePolicy.');
        }

        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling putCertificatePolicy.');
        }

        if (requestParameters.createCertificatePolicyRequest === null || requestParameters.createCertificatePolicyRequest === undefined) {
            throw new runtime.RequiredError('createCertificatePolicyRequest','Required parameter requestParameters.createCertificatePolicyRequest was null or undefined when calling putCertificatePolicy.');
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
            path: `/v2/{namespaceProvider}/{namespaceId}/certificate-policies/{id}`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: CreateCertificatePolicyRequestToJSON(requestParameters.createCertificatePolicyRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CertificatePolicyFromJSON(jsonValue));
    }

    /**
     * put certificate policy
     */
    async putCertificatePolicy(requestParameters: PutCertificatePolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CertificatePolicy> {
        const response = await this.putCertificatePolicyRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * put key policy
     */
    async putKeyPolicyRaw(requestParameters: PutKeyPolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<KeyPolicy>> {
        if (requestParameters.namespaceProvider === null || requestParameters.namespaceProvider === undefined) {
            throw new runtime.RequiredError('namespaceProvider','Required parameter requestParameters.namespaceProvider was null or undefined when calling putKeyPolicy.');
        }

        if (requestParameters.namespaceId === null || requestParameters.namespaceId === undefined) {
            throw new runtime.RequiredError('namespaceId','Required parameter requestParameters.namespaceId was null or undefined when calling putKeyPolicy.');
        }

        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling putKeyPolicy.');
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
            path: `/v2/{namespaceProvider}/{namespaceId}/key-policies/{id}`.replace(`{${"namespaceProvider"}}`, encodeURIComponent(String(requestParameters.namespaceProvider))).replace(`{${"namespaceId"}}`, encodeURIComponent(String(requestParameters.namespaceId))).replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => KeyPolicyFromJSON(jsonValue));
    }

    /**
     * put key policy
     */
    async putKeyPolicy(requestParameters: PutKeyPolicyRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<KeyPolicy> {
        const response = await this.putKeyPolicyRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Sync managed app
     */
    async syncSystemAppRaw(requestParameters: SyncSystemAppRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ApplicationByAppId>> {
        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling syncSystemApp.');
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
            path: `/v2/system-apps/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ApplicationByAppIdFromJSON(jsonValue));
    }

    /**
     * Sync managed app
     */
    async syncSystemApp(requestParameters: SyncSystemAppRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ApplicationByAppId> {
        const response = await this.syncSystemAppRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
