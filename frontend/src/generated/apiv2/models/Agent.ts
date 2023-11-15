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
/**
 * 
 * @export
 * @interface Agent
 */
export interface Agent {
    /**
     * 
     * @type {string}
     * @memberof Agent
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof Agent
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof Agent
     */
    updatedBy: string;
    /**
     * 
     * @type {Date}
     * @memberof Agent
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof Agent
     */
    displayName?: string;
    /**
     * 
     * @type {string}
     * @memberof Agent
     */
    applicationId: string;
    /**
     * 
     * @type {string}
     * @memberof Agent
     */
    servicePrincipalId: string;
}

/**
 * Check if a given object implements the Agent interface.
 */
export function instanceOfAgent(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "applicationId" in value;
    isInstance = isInstance && "servicePrincipalId" in value;

    return isInstance;
}

export function AgentFromJSON(json: any): Agent {
    return AgentFromJSONTyped(json, false);
}

export function AgentFromJSONTyped(json: any, ignoreDiscriminator: boolean): Agent {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'applicationId': json['applicationId'],
        'servicePrincipalId': json['servicePrincipalId'],
    };
}

export function AgentToJSON(value?: Agent | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'updated': (value.updated.toISOString()),
        'updatedBy': value.updatedBy,
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'displayName': value.displayName,
        'applicationId': value.applicationId,
        'servicePrincipalId': value.servicePrincipalId,
    };
}

