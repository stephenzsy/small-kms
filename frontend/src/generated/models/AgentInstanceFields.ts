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
 * @interface AgentInstanceFields
 */
export interface AgentInstanceFields {
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceFields
     */
    endpoint?: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceFields
     */
    version: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceFields
     */
    buildId: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceFields
     */
    mode?: string;
}

/**
 * Check if a given object implements the AgentInstanceFields interface.
 */
export function instanceOfAgentInstanceFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "version" in value;
    isInstance = isInstance && "buildId" in value;

    return isInstance;
}

export function AgentInstanceFieldsFromJSON(json: any): AgentInstanceFields {
    return AgentInstanceFieldsFromJSONTyped(json, false);
}

export function AgentInstanceFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentInstanceFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'endpoint': !exists(json, 'endpoint') ? undefined : json['endpoint'],
        'version': json['version'],
        'buildId': json['buildId'],
        'mode': !exists(json, 'mode') ? undefined : json['mode'],
    };
}

export function AgentInstanceFieldsToJSON(value?: AgentInstanceFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'endpoint': value.endpoint,
        'version': value.version,
        'buildId': value.buildId,
        'mode': value.mode,
    };
}

