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
/**
 * 
 * @export
 * @interface AgentCheckInResult
 */
export interface AgentCheckInResult {
    /**
     * 
     * @type {string}
     * @memberof AgentCheckInResult
     */
    message?: string;
}

/**
 * Check if a given object implements the AgentCheckInResult interface.
 */
export function instanceOfAgentCheckInResult(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function AgentCheckInResultFromJSON(json: any): AgentCheckInResult {
    return AgentCheckInResultFromJSONTyped(json, false);
}

export function AgentCheckInResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentCheckInResult {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'message': !exists(json, 'message') ? undefined : json['message'],
    };
}

export function AgentCheckInResultToJSON(value?: AgentCheckInResult | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'message': value.message,
    };
}

