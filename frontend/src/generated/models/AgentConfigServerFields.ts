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
 * @interface AgentConfigServerFields
 */
export interface AgentConfigServerFields {
    /**
     * 
     * @type {string}
     * @memberof AgentConfigServerFields
     */
    tlsCertificateId: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof AgentConfigServerFields
     */
    jwtKeyCertIds: Array<string>;
    /**
     * 
     * @type {string}
     * @memberof AgentConfigServerFields
     */
    imageTag: string;
}

/**
 * Check if a given object implements the AgentConfigServerFields interface.
 */
export function instanceOfAgentConfigServerFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "tlsCertificateId" in value;
    isInstance = isInstance && "jwtKeyCertIds" in value;
    isInstance = isInstance && "imageTag" in value;

    return isInstance;
}

export function AgentConfigServerFieldsFromJSON(json: any): AgentConfigServerFields {
    return AgentConfigServerFieldsFromJSONTyped(json, false);
}

export function AgentConfigServerFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigServerFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'tlsCertificateId': json['tlsCertificateId'],
        'jwtKeyCertIds': json['jwtKeyCertIds'],
        'imageTag': json['imageTag'],
    };
}

export function AgentConfigServerFieldsToJSON(value?: AgentConfigServerFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'tlsCertificateId': value.tlsCertificateId,
        'jwtKeyCertIds': value.jwtKeyCertIds,
        'imageTag': value.imageTag,
    };
}
