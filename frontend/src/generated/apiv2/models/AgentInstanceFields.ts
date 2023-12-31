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
 * @interface AgentInstanceFields
 */
export interface AgentInstanceFields {
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceFields
     */
    tlsCertificateId: string;
    /**
     * 
     * @type {boolean}
     * @memberof AgentInstanceFields
     */
    tlsCertificateSignedByPublicCa: boolean;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceFields
     */
    jwtVerifyKeyId: string;
}

/**
 * Check if a given object implements the AgentInstanceFields interface.
 */
export function instanceOfAgentInstanceFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "tlsCertificateId" in value;
    isInstance = isInstance && "tlsCertificateSignedByPublicCa" in value;
    isInstance = isInstance && "jwtVerifyKeyId" in value;

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
        
        'tlsCertificateId': json['tlsCertificateId'],
        'tlsCertificateSignedByPublicCa': json['tlsCertificateSignedByPublicCa'],
        'jwtVerifyKeyId': json['jwtVerifyKeyId'],
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
        
        'tlsCertificateId': value.tlsCertificateId,
        'tlsCertificateSignedByPublicCa': value.tlsCertificateSignedByPublicCa,
        'jwtVerifyKeyId': value.jwtVerifyKeyId,
    };
}

