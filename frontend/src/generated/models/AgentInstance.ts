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
import type { AgentMode } from './AgentMode';
import {
    AgentModeFromJSON,
    AgentModeFromJSONTyped,
    AgentModeToJSON,
} from './AgentMode';

/**
 * 
 * @export
 * @interface AgentInstance
 */
export interface AgentInstance {
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof AgentInstance
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof AgentInstance
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    updatedBy?: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    endpoint?: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    version: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    buildId: string;
    /**
     * 
     * @type {AgentMode}
     * @memberof AgentInstance
     */
    mode: AgentMode;
}

/**
 * Check if a given object implements the AgentInstance interface.
 */
export function instanceOfAgentInstance(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "version" in value;
    isInstance = isInstance && "buildId" in value;
    isInstance = isInstance && "mode" in value;

    return isInstance;
}

export function AgentInstanceFromJSON(json: any): AgentInstance {
    return AgentInstanceFromJSONTyped(json, false);
}

export function AgentInstanceFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentInstance {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'endpoint': !exists(json, 'endpoint') ? undefined : json['endpoint'],
        'version': json['version'],
        'buildId': json['buildId'],
        'mode': AgentModeFromJSON(json['mode']),
    };
}

export function AgentInstanceToJSON(value?: AgentInstance | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
        'endpoint': value.endpoint,
        'version': value.version,
        'buildId': value.buildId,
        'mode': AgentModeToJSON(value.mode),
    };
}

