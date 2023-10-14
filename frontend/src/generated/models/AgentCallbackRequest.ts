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

import { exists, mapValues } from '../runtime';
import type { ServiceRuntimeInfo } from './ServiceRuntimeInfo';
import {
    ServiceRuntimeInfoFromJSON,
    ServiceRuntimeInfoFromJSONTyped,
    ServiceRuntimeInfoToJSON,
} from './ServiceRuntimeInfo';

/**
 * 
 * @export
 * @interface AgentCallbackRequest
 */
export interface AgentCallbackRequest {
    /**
     * 
     * @type {string}
     * @memberof AgentCallbackRequest
     */
    configVersion?: string;
    /**
     * 
     * @type {ServiceRuntimeInfo}
     * @memberof AgentCallbackRequest
     */
    serviceRuntime?: ServiceRuntimeInfo;
}

/**
 * Check if a given object implements the AgentCallbackRequest interface.
 */
export function instanceOfAgentCallbackRequest(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function AgentCallbackRequestFromJSON(json: any): AgentCallbackRequest {
    return AgentCallbackRequestFromJSONTyped(json, false);
}

export function AgentCallbackRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentCallbackRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'configVersion': !exists(json, 'configVersion') ? undefined : json['configVersion'],
        'serviceRuntime': !exists(json, 'serviceRuntime') ? undefined : ServiceRuntimeInfoFromJSON(json['serviceRuntime']),
    };
}

export function AgentCallbackRequestToJSON(value?: AgentCallbackRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'configVersion': value.configVersion,
        'serviceRuntime': ServiceRuntimeInfoToJSON(value.serviceRuntime),
    };
}

