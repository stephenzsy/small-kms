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
  Ref,
} from '../models/index';
import {
    RefFromJSON,
    RefToJSON,
} from '../models/index';

/**
 * 
 */
export class AdminApi extends runtime.BaseAPI {

    /**
     * List agents
     */
    async listAgentsRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<Ref>>> {
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
            path: `/v2/agents`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RefFromJSON));
    }

    /**
     * List agents
     */
    async listAgents(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<Ref>> {
        const response = await this.listAgentsRaw(initOverrides);
        return await response.value();
    }

}
