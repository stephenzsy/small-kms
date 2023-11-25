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
 * @interface AgentConfigEndpointFields
 */
export interface AgentConfigEndpointFields {
    /**
     * 
     * @type {string}
     * @memberof AgentConfigEndpointFields
     */
    tlsCertificatePolicyId: string;
    /**
     * 
     * @type {string}
     * @memberof AgentConfigEndpointFields
     */
    jwtVerifyKeyPolicyId: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof AgentConfigEndpointFields
     */
    jwtVerifyKeyIds?: Array<string>;
}

/**
 * Check if a given object implements the AgentConfigEndpointFields interface.
 */
export function instanceOfAgentConfigEndpointFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "tlsCertificatePolicyId" in value;
    isInstance = isInstance && "jwtVerifyKeyPolicyId" in value;

    return isInstance;
}

export function AgentConfigEndpointFieldsFromJSON(json: any): AgentConfigEndpointFields {
    return AgentConfigEndpointFieldsFromJSONTyped(json, false);
}

export function AgentConfigEndpointFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigEndpointFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'tlsCertificatePolicyId': json['tlsCertificatePolicyId'],
        'jwtVerifyKeyPolicyId': json['jwtVerifyKeyPolicyId'],
        'jwtVerifyKeyIds': !exists(json, 'jwtVerifyKeyIds') ? undefined : json['jwtVerifyKeyIds'],
    };
}

export function AgentConfigEndpointFieldsToJSON(value?: AgentConfigEndpointFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'tlsCertificatePolicyId': value.tlsCertificatePolicyId,
        'jwtVerifyKeyPolicyId': value.jwtVerifyKeyPolicyId,
        'jwtVerifyKeyIds': value.jwtVerifyKeyIds,
    };
}
