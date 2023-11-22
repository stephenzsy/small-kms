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

import { exists, mapValues } from '../runtime';
import type { AgentConfigName } from './AgentConfigName';
import {
    AgentConfigNameFromJSON,
    AgentConfigNameFromJSONTyped,
    AgentConfigNameToJSON,
} from './AgentConfigName';

/**
 * 
 * @export
 * @interface AgentConfigRef
 */
export interface AgentConfigRef {
    /**
     * 
     * @type {AgentConfigName}
     * @memberof AgentConfigRef
     */
    name: AgentConfigName;
    /**
     * 
     * @type {Date}
     * @memberof AgentConfigRef
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof AgentConfigRef
     */
    version: string;
}

/**
 * Check if a given object implements the AgentConfigRef interface.
 */
export function instanceOfAgentConfigRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "version" in value;

    return isInstance;
}

export function AgentConfigRefFromJSON(json: any): AgentConfigRef {
    return AgentConfigRefFromJSONTyped(json, false);
}

export function AgentConfigRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'name': AgentConfigNameFromJSON(json['name']),
        'updated': (new Date(json['updated'])),
        'version': json['version'],
    };
}

export function AgentConfigRefToJSON(value?: AgentConfigRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'name': AgentConfigNameToJSON(value.name),
        'updated': (value.updated.toISOString()),
        'version': value.version,
    };
}

