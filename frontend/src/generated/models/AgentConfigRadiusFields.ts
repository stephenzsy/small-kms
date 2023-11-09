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
import type { RadiusClientConfig } from './RadiusClientConfig';
import {
    RadiusClientConfigFromJSON,
    RadiusClientConfigFromJSONTyped,
    RadiusClientConfigToJSON,
} from './RadiusClientConfig';

/**
 * 
 * @export
 * @interface AgentConfigRadiusFields
 */
export interface AgentConfigRadiusFields {
    /**
     * 
     * @type {string}
     * @memberof AgentConfigRadiusFields
     */
    azureAcrImageRef?: string;
    /**
     * 
     * @type {Array<RadiusClientConfig>}
     * @memberof AgentConfigRadiusFields
     */
    clients?: Array<RadiusClientConfig>;
}

/**
 * Check if a given object implements the AgentConfigRadiusFields interface.
 */
export function instanceOfAgentConfigRadiusFields(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function AgentConfigRadiusFieldsFromJSON(json: any): AgentConfigRadiusFields {
    return AgentConfigRadiusFieldsFromJSONTyped(json, false);
}

export function AgentConfigRadiusFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigRadiusFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'azureAcrImageRef': !exists(json, 'azureAcrImageRef') ? undefined : json['azureAcrImageRef'],
        'clients': !exists(json, 'clients') ? undefined : ((json['clients'] as Array<any>).map(RadiusClientConfigFromJSON)),
    };
}

export function AgentConfigRadiusFieldsToJSON(value?: AgentConfigRadiusFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'azureAcrImageRef': value.azureAcrImageRef,
        'clients': value.clients === undefined ? undefined : ((value.clients as Array<any>).map(RadiusClientConfigToJSON)),
    };
}
