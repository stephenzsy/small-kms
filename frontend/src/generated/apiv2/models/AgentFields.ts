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
 * @interface AgentFields
 */
export interface AgentFields {
    /**
     * 
     * @type {string}
     * @memberof AgentFields
     */
    clientCredentialCertificatePolicyId: string;
}

/**
 * Check if a given object implements the AgentFields interface.
 */
export function instanceOfAgentFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "clientCredentialCertificatePolicyId" in value;

    return isInstance;
}

export function AgentFieldsFromJSON(json: any): AgentFields {
    return AgentFieldsFromJSONTyped(json, false);
}

export function AgentFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'clientCredentialCertificatePolicyId': json['clientCredentialCertificatePolicyId'],
    };
}

export function AgentFieldsToJSON(value?: AgentFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'clientCredentialCertificatePolicyId': value.clientCredentialCertificatePolicyId,
    };
}
