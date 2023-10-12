/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  AgentCheckInResult,
  AgentConfigName,
  AgentConfiguration,
  AgentHostRole,
  ServiceConfig,
} from '../models';
import {
    AgentCheckInResultFromJSON,
    AgentCheckInResultToJSON,
    AgentConfigNameFromJSON,
    AgentConfigNameToJSON,
    AgentConfigurationFromJSON,
    AgentConfigurationToJSON,
    AgentHostRoleFromJSON,
    AgentHostRoleToJSON,
    ServiceConfigFromJSON,
    ServiceConfigToJSON,
} from '../models';

export interface AgentCheckInRequest {
    hostRoles?: Array<AgentHostRole>;
}

export interface AgentGetConfigurationRequest {
    configName: AgentConfigName;
    xSmallkmsIfVersionNotMatch?: string;
    refreshToken?: string;
}

/**
 * 
 */
export class AgentApi extends runtime.BaseAPI {

    /**
     */
    async agentCheckInRaw(requestParameters: AgentCheckInRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AgentCheckInResult>> {
        const queryParameters: any = {};

        if (requestParameters.hostRoles) {
            queryParameters['hostRoles'] = requestParameters.hostRoles;
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
            path: `/v3/agent/check-in`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AgentCheckInResultFromJSON(jsonValue));
    }

    /**
     */
    async agentCheckIn(requestParameters: AgentCheckInRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AgentCheckInResult> {
        const response = await this.agentCheckInRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get agent configuration
     */
    async agentGetConfigurationRaw(requestParameters: AgentGetConfigurationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AgentConfiguration>> {
        if (requestParameters.configName === null || requestParameters.configName === undefined) {
            throw new runtime.RequiredError('configName','Required parameter requestParameters.configName was null or undefined when calling agentGetConfiguration.');
        }

        const queryParameters: any = {};

        if (requestParameters.refreshToken !== undefined) {
            queryParameters['refreshToken'] = requestParameters.refreshToken;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (requestParameters.xSmallkmsIfVersionNotMatch !== undefined && requestParameters.xSmallkmsIfVersionNotMatch !== null) {
            headerParameters['X-Smallkms-If-Version-Not-Match'] = String(requestParameters.xSmallkmsIfVersionNotMatch);
        }

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v3/agent/config/{configName}`.replace(`{${"configName"}}`, encodeURIComponent(String(requestParameters.configName))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AgentConfigurationFromJSON(jsonValue));
    }

    /**
     * Get agent configuration
     */
    async agentGetConfiguration(requestParameters: AgentGetConfigurationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AgentConfiguration> {
        const response = await this.agentGetConfigurationRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get service config
     */
    async getServiceConfigRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ServiceConfig>> {
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
            path: `/v3/service/config`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ServiceConfigFromJSON(jsonValue));
    }

    /**
     * Get service config
     */
    async getServiceConfig(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ServiceConfig> {
        const response = await this.getServiceConfigRaw(initOverrides);
        return await response.value();
    }

}