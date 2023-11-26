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
import type { AgentInstanceState } from './AgentInstanceState';
import {
    AgentInstanceStateFromJSON,
    AgentInstanceStateFromJSONTyped,
    AgentInstanceStateToJSON,
} from './AgentInstanceState';

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
     * @type {string}
     * @memberof AgentInstance
     */
    updatedBy: string;
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
    displayName?: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    endpoint: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    buildId: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    configVersion: string;
    /**
     * 
     * @type {AgentInstanceState}
     * @memberof AgentInstance
     */
    state: AgentInstanceState;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    tlsCertificateId: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstance
     */
    jwtVerifyKeyId: string;
}

/**
 * Check if a given object implements the AgentInstance interface.
 */
export function instanceOfAgentInstance(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "endpoint" in value;
    isInstance = isInstance && "buildId" in value;
    isInstance = isInstance && "configVersion" in value;
    isInstance = isInstance && "state" in value;
    isInstance = isInstance && "tlsCertificateId" in value;
    isInstance = isInstance && "jwtVerifyKeyId" in value;

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
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'endpoint': json['endpoint'],
        'buildId': json['buildId'],
        'configVersion': json['configVersion'],
        'state': AgentInstanceStateFromJSON(json['state']),
        'tlsCertificateId': json['tlsCertificateId'],
        'jwtVerifyKeyId': json['jwtVerifyKeyId'],
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
        'updatedBy': value.updatedBy,
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'displayName': value.displayName,
        'endpoint': value.endpoint,
        'buildId': value.buildId,
        'configVersion': value.configVersion,
        'state': AgentInstanceStateToJSON(value.state),
        'tlsCertificateId': value.tlsCertificateId,
        'jwtVerifyKeyId': value.jwtVerifyKeyId,
    };
}

