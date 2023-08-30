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
import type { KeyParameters } from './KeyParameters';
import {
    KeyParametersFromJSON,
    KeyParametersFromJSONTyped,
    KeyParametersToJSON,
} from './KeyParameters';
import type { PolicyId } from './PolicyId';
import {
    PolicyIdFromJSON,
    PolicyIdFromJSONTyped,
    PolicyIdToJSON,
} from './PolicyId';

/**
 * 
 * @export
 * @interface CertificateEnrollmentPolicy
 */
export interface CertificateEnrollmentPolicy {
    /**
     * Unique ID of the namespace
     * @type {string}
     * @memberof CertificateEnrollmentPolicy
     */
    namespaceId: string;
    /**
     * 
     * @type {PolicyId}
     * @memberof CertificateEnrollmentPolicy
     */
    id: PolicyId;
    /**
     * Unique ID of the user who created the policy
     * @type {string}
     * @memberof CertificateEnrollmentPolicy
     */
    updatedBy: string;
    /**
     * Time when the policy was last updated
     * @type {Date}
     * @memberof CertificateEnrollmentPolicy
     */
    updatedAt: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificateEnrollmentPolicy
     */
    issuerNamespace: string;
    /**
     * 
     * @type {string}
     * @memberof CertificateEnrollmentPolicy
     */
    issuerId: string;
    /**
     * RFC3339 duration string
     * @type {string}
     * @memberof CertificateEnrollmentPolicy
     */
    validity: string;
    /**
     * 
     * @type {KeyParameters}
     * @memberof CertificateEnrollmentPolicy
     */
    keyParameters: KeyParameters;
    /**
     * 
     * @type {string}
     * @memberof CertificateEnrollmentPolicy
     */
    delegatedService?: string;
}

/**
 * Check if a given object implements the CertificateEnrollmentPolicy interface.
 */
export function instanceOfCertificateEnrollmentPolicy(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceId" in value;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "updatedAt" in value;
    isInstance = isInstance && "issuerNamespace" in value;
    isInstance = isInstance && "issuerId" in value;
    isInstance = isInstance && "validity" in value;
    isInstance = isInstance && "keyParameters" in value;

    return isInstance;
}

export function CertificateEnrollmentPolicyFromJSON(json: any): CertificateEnrollmentPolicy {
    return CertificateEnrollmentPolicyFromJSONTyped(json, false);
}

export function CertificateEnrollmentPolicyFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateEnrollmentPolicy {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'namespaceId': json['namespaceId'],
        'id': PolicyIdFromJSON(json['id']),
        'updatedBy': json['updatedBy'],
        'updatedAt': (new Date(json['updatedAt'])),
        'issuerNamespace': json['issuerNamespace'],
        'issuerId': json['issuerId'],
        'validity': json['validity'],
        'keyParameters': KeyParametersFromJSON(json['keyParameters']),
        'delegatedService': !exists(json, 'delegatedService') ? undefined : json['delegatedService'],
    };
}

export function CertificateEnrollmentPolicyToJSON(value?: CertificateEnrollmentPolicy | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'namespaceId': value.namespaceId,
        'id': PolicyIdToJSON(value.id),
        'updatedBy': value.updatedBy,
        'updatedAt': (value.updatedAt.toISOString()),
        'issuerNamespace': value.issuerNamespace,
        'issuerId': value.issuerId,
        'validity': value.validity,
        'keyParameters': KeyParametersToJSON(value.keyParameters),
        'delegatedService': value.delegatedService,
    };
}

