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
import type { CertificateFlag } from './CertificateFlag';
import {
    CertificateFlagFromJSON,
    CertificateFlagFromJSONTyped,
    CertificateFlagToJSON,
} from './CertificateFlag';
import type { CertificatePrivateKeyMode } from './CertificatePrivateKeyMode';
import {
    CertificatePrivateKeyModeFromJSON,
    CertificatePrivateKeyModeFromJSONTyped,
    CertificatePrivateKeyModeToJSON,
} from './CertificatePrivateKeyMode';
import type { CertificateSubject } from './CertificateSubject';
import {
    CertificateSubjectFromJSON,
    CertificateSubjectFromJSONTyped,
    CertificateSubjectToJSON,
} from './CertificateSubject';
import type { JsonWebKeySpec } from './JsonWebKeySpec';
import {
    JsonWebKeySpecFromJSON,
    JsonWebKeySpecFromJSONTyped,
    JsonWebKeySpecToJSON,
} from './JsonWebKeySpec';
import type { SubjectAlternativeNames } from './SubjectAlternativeNames';
import {
    SubjectAlternativeNamesFromJSON,
    SubjectAlternativeNamesFromJSONTyped,
    SubjectAlternativeNamesToJSON,
} from './SubjectAlternativeNames';

/**
 * 
 * @export
 * @interface CreateCertificatePolicyRequest
 */
export interface CreateCertificatePolicyRequest {
    /**
     * 
     * @type {string}
     * @memberof CreateCertificatePolicyRequest
     */
    displayName?: string;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificatePolicyRequest
     */
    issuerPolicyIdentifier?: string;
    /**
     * 
     * @type {JsonWebKeySpec}
     * @memberof CreateCertificatePolicyRequest
     */
    keySpec?: JsonWebKeySpec;
    /**
     * 
     * @type {CertificatePrivateKeyMode}
     * @memberof CreateCertificatePolicyRequest
     */
    keyMode?: CertificatePrivateKeyMode;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificatePolicyRequest
     */
    expiryTime?: string;
    /**
     * 
     * @type {CertificateSubject}
     * @memberof CreateCertificatePolicyRequest
     */
    subject: CertificateSubject;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof CreateCertificatePolicyRequest
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateFlag>}
     * @memberof CreateCertificatePolicyRequest
     */
    flags?: Array<CertificateFlag>;
}

/**
 * Check if a given object implements the CreateCertificatePolicyRequest interface.
 */
export function instanceOfCreateCertificatePolicyRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "subject" in value;

    return isInstance;
}

export function CreateCertificatePolicyRequestFromJSON(json: any): CreateCertificatePolicyRequest {
    return CreateCertificatePolicyRequestFromJSONTyped(json, false);
}

export function CreateCertificatePolicyRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateCertificatePolicyRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'issuerPolicyIdentifier': !exists(json, 'issuerPolicyIdentifier') ? undefined : json['issuerPolicyIdentifier'],
        'keySpec': !exists(json, 'keySpec') ? undefined : JsonWebKeySpecFromJSON(json['keySpec']),
        'keyMode': !exists(json, 'keyMode') ? undefined : CertificatePrivateKeyModeFromJSON(json['keyMode']),
        'expiryTime': !exists(json, 'expiryTime') ? undefined : json['expiryTime'],
        'subject': CertificateSubjectFromJSON(json['subject']),
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : SubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'flags': !exists(json, 'flags') ? undefined : ((json['flags'] as Array<any>).map(CertificateFlagFromJSON)),
    };
}

export function CreateCertificatePolicyRequestToJSON(value?: CreateCertificatePolicyRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'displayName': value.displayName,
        'issuerPolicyIdentifier': value.issuerPolicyIdentifier,
        'keySpec': JsonWebKeySpecToJSON(value.keySpec),
        'keyMode': CertificatePrivateKeyModeToJSON(value.keyMode),
        'expiryTime': value.expiryTime,
        'subject': CertificateSubjectToJSON(value.subject),
        'subjectAlternativeNames': SubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'flags': value.flags === undefined ? undefined : ((value.flags as Array<any>).map(CertificateFlagToJSON)),
    };
}

