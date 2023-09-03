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
import type { CertificateSubject } from './CertificateSubject';
import {
    CertificateSubjectFromJSON,
    CertificateSubjectFromJSONTyped,
    CertificateSubjectToJSON,
} from './CertificateSubject';
import type { CertificateSubjectAlternativeNames } from './CertificateSubjectAlternativeNames';
import {
    CertificateSubjectAlternativeNamesFromJSON,
    CertificateSubjectAlternativeNamesFromJSONTyped,
    CertificateSubjectAlternativeNamesToJSON,
} from './CertificateSubjectAlternativeNames';
import type { CertificateUsage } from './CertificateUsage';
import {
    CertificateUsageFromJSON,
    CertificateUsageFromJSONTyped,
    CertificateUsageToJSON,
} from './CertificateUsage';
import type { KeyProperties } from './KeyProperties';
import {
    KeyPropertiesFromJSON,
    KeyPropertiesFromJSONTyped,
    KeyPropertiesToJSON,
} from './KeyProperties';

/**
 * 
 * @export
 * @interface CertificateRequestPolicyParameters
 */
export interface CertificateRequestPolicyParameters {
    /**
     * ID of the issuer namespace
     * @type {string}
     * @memberof CertificateRequestPolicyParameters
     */
    issuerNamespaceId: string;
    /**
     * 
     * @type {number}
     * @memberof CertificateRequestPolicyParameters
     */
    validityMonths?: number;
    /**
     * 
     * @type {KeyProperties}
     * @memberof CertificateRequestPolicyParameters
     */
    keyProperties?: KeyProperties;
    /**
     * 
     * @type {CertificateSubject}
     * @memberof CertificateRequestPolicyParameters
     */
    subject: CertificateSubject;
    /**
     * 
     * @type {CertificateSubjectAlternativeNames}
     * @memberof CertificateRequestPolicyParameters
     */
    subjectAlternativeNames?: CertificateSubjectAlternativeNames;
    /**
     * 
     * @type {CertificateUsage}
     * @memberof CertificateRequestPolicyParameters
     */
    usage: CertificateUsage;
    /**
     * Number of days left to trigger renewal, (0, 1) range indicates percentage
     * @type {number}
     * @memberof CertificateRequestPolicyParameters
     */
    autoRenewalThreshold?: number;
    /**
     * 
     * @type {string}
     * @memberof CertificateRequestPolicyParameters
     */
    keyStorePath: string;
}

/**
 * Check if a given object implements the CertificateRequestPolicyParameters interface.
 */
export function instanceOfCertificateRequestPolicyParameters(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "issuerNamespaceId" in value;
    isInstance = isInstance && "subject" in value;
    isInstance = isInstance && "usage" in value;
    isInstance = isInstance && "keyStorePath" in value;

    return isInstance;
}

export function CertificateRequestPolicyParametersFromJSON(json: any): CertificateRequestPolicyParameters {
    return CertificateRequestPolicyParametersFromJSONTyped(json, false);
}

export function CertificateRequestPolicyParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateRequestPolicyParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'issuerNamespaceId': json['issuerNamespaceId'],
        'validityMonths': !exists(json, 'validity_months') ? undefined : json['validity_months'],
        'keyProperties': !exists(json, 'keyProperties') ? undefined : KeyPropertiesFromJSON(json['keyProperties']),
        'subject': CertificateSubjectFromJSON(json['subject']),
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : CertificateSubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'usage': CertificateUsageFromJSON(json['usage']),
        'autoRenewalThreshold': !exists(json, 'autoRenewalThreshold') ? undefined : json['autoRenewalThreshold'],
        'keyStorePath': json['keyStorePath'],
    };
}

export function CertificateRequestPolicyParametersToJSON(value?: CertificateRequestPolicyParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'issuerNamespaceId': value.issuerNamespaceId,
        'validity_months': value.validityMonths,
        'keyProperties': KeyPropertiesToJSON(value.keyProperties),
        'subject': CertificateSubjectToJSON(value.subject),
        'subjectAlternativeNames': CertificateSubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'usage': CertificateUsageToJSON(value.usage),
        'autoRenewalThreshold': value.autoRenewalThreshold,
        'keyStorePath': value.keyStorePath,
    };
}

