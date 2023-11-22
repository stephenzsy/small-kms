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
 * @interface AgentConfigIdentityFields
 */
export interface AgentConfigIdentityFields {
    /**
     * 
     * @type {string}
     * @memberof AgentConfigIdentityFields
     */
    keyCredentialCertificatePolicyId: string;
}

/**
 * Check if a given object implements the AgentConfigIdentityFields interface.
 */
export function instanceOfAgentConfigIdentityFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "keyCredentialCertificatePolicyId" in value;

    return isInstance;
}

export function AgentConfigIdentityFieldsFromJSON(json: any): AgentConfigIdentityFields {
    return AgentConfigIdentityFieldsFromJSONTyped(json, false);
}

export function AgentConfigIdentityFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigIdentityFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'keyCredentialCertificatePolicyId': json['keyCredentialCertificatePolicyId'],
    };
}

export function AgentConfigIdentityFieldsToJSON(value?: AgentConfigIdentityFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'keyCredentialCertificatePolicyId': value.keyCredentialCertificatePolicyId,
    };
}

