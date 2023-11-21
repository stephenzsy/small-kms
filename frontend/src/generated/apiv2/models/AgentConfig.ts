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
 * @interface AgentConfig
 */
export interface AgentConfig {
    /**
     * 
     * @type {string}
     * @memberof AgentConfig
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof AgentConfig
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof AgentConfig
     */
    updatedBy: string;
    /**
     * 
     * @type {Date}
     * @memberof AgentConfig
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof AgentConfig
     */
    displayName?: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof AgentConfig
     */
    envGuards: Array<string>;
    /**
     * 
     * @type {string}
     * @memberof AgentConfig
     */
    keyCredentialsCertificatePolicyId: string;
}

/**
 * Check if a given object implements the AgentConfig interface.
 */
export function instanceOfAgentConfig(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "envGuards" in value;
    isInstance = isInstance && "keyCredentialsCertificatePolicyId" in value;

    return isInstance;
}

export function AgentConfigFromJSON(json: any): AgentConfig {
    return AgentConfigFromJSONTyped(json, false);
}

export function AgentConfigFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfig {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'envGuards': json['envGuards'],
        'keyCredentialsCertificatePolicyId': json['keyCredentialsCertificatePolicyId'],
    };
}

export function AgentConfigToJSON(value?: AgentConfig | null): any {
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
        'envGuards': value.envGuards,
        'keyCredentialsCertificatePolicyId': value.keyCredentialsCertificatePolicyId,
    };
}

