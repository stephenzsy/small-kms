/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { CertificateRequestPolicyParameters } from './CertificateRequestPolicyParameters';
import {
    CertificateRequestPolicyParametersFromJSON,
    CertificateRequestPolicyParametersFromJSONTyped,
    CertificateRequestPolicyParametersToJSON,
} from './CertificateRequestPolicyParameters';
import type { PolicyType } from './PolicyType';
import {
    PolicyTypeFromJSON,
    PolicyTypeFromJSONTyped,
    PolicyTypeToJSON,
} from './PolicyType';

/**
 * 
 * @export
 * @interface Policy
 */
export interface Policy {
    /**
     * Unique ID of the namespace
     * @type {string}
     * @memberof Policy
     */
    namespaceId: string;
    /**
     * 
     * @type {string}
     * @memberof Policy
     */
    id: string;
    /**
     * Unique ID of the user who created the policy
     * @type {string}
     * @memberof Policy
     */
    updatedBy: string;
    /**
     * Time when the policy was last updated
     * @type {Date}
     * @memberof Policy
     */
    updated: Date;
    /**
     * 
     * @type {PolicyType}
     * @memberof Policy
     */
    policyType: PolicyType;
    /**
     * 
     * @type {CertificateRequestPolicyParameters}
     * @memberof Policy
     */
    certRequest?: CertificateRequestPolicyParameters;
}

/**
 * Check if a given object implements the Policy interface.
 */
export function instanceOfPolicy(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceId" in value;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "policyType" in value;

    return isInstance;
}

export function PolicyFromJSON(json: any): Policy {
    return PolicyFromJSONTyped(json, false);
}

export function PolicyFromJSONTyped(json: any, ignoreDiscriminator: boolean): Policy {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'namespaceId': json['namespaceId'],
        'id': json['id'],
        'updatedBy': json['updatedBy'],
        'updated': (new Date(json['updated'])),
        'policyType': PolicyTypeFromJSON(json['policyType']),
        'certRequest': !exists(json, 'certRequest') ? undefined : CertificateRequestPolicyParametersFromJSON(json['certRequest']),
    };
}

export function PolicyToJSON(value?: Policy | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'namespaceId': value.namespaceId,
        'id': value.id,
        'updatedBy': value.updatedBy,
        'updated': (value.updated.toISOString()),
        'policyType': PolicyTypeToJSON(value.policyType),
        'certRequest': CertificateRequestPolicyParametersToJSON(value.certRequest),
    };
}

