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
import type { JsonWebSignatureAlgorithm } from './JsonWebSignatureAlgorithm';
import {
    JsonWebSignatureAlgorithmFromJSON,
    JsonWebSignatureAlgorithmFromJSONTyped,
    JsonWebSignatureAlgorithmToJSON,
} from './JsonWebSignatureAlgorithm';
import type { SubjectAlternativeNames } from './SubjectAlternativeNames';
import {
    SubjectAlternativeNamesFromJSON,
    SubjectAlternativeNamesFromJSONTyped,
    SubjectAlternativeNamesToJSON,
} from './SubjectAlternativeNames';

/**
 * 
 * @export
 * @interface CertificatePolicy
 */
export interface CertificatePolicy {
    /**
     * 
     * @type {string}
     * @memberof CertificatePolicy
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof CertificatePolicy
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificatePolicy
     */
    updatedBy: string;
    /**
     * 
     * @type {Date}
     * @memberof CertificatePolicy
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificatePolicy
     */
    displayName?: string;
    /**
     * 
     * @type {JsonWebKeySpec}
     * @memberof CertificatePolicy
     */
    keySpec: JsonWebKeySpec;
    /**
     * 
     * @type {boolean}
     * @memberof CertificatePolicy
     */
    ext: boolean;
    /**
     * 
     * @type {string}
     * @memberof CertificatePolicy
     */
    expiryTime: string;
    /**
     * 
     * @type {JsonWebSignatureAlgorithm}
     * @memberof CertificatePolicy
     */
    alg: JsonWebSignatureAlgorithm;
    /**
     * 
     * @type {string}
     * @memberof CertificatePolicy
     */
    issuerIdentifier: string;
    /**
     * 
     * @type {CertificateSubject}
     * @memberof CertificatePolicy
     */
    subject: CertificateSubject;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof CertificatePolicy
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateFlag>}
     * @memberof CertificatePolicy
     */
    flags: Array<CertificateFlag>;
}

/**
 * Check if a given object implements the CertificatePolicy interface.
 */
export function instanceOfCertificatePolicy(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "keySpec" in value;
    isInstance = isInstance && "ext" in value;
    isInstance = isInstance && "expiryTime" in value;
    isInstance = isInstance && "alg" in value;
    isInstance = isInstance && "issuerIdentifier" in value;
    isInstance = isInstance && "subject" in value;
    isInstance = isInstance && "flags" in value;

    return isInstance;
}

export function CertificatePolicyFromJSON(json: any): CertificatePolicy {
    return CertificatePolicyFromJSONTyped(json, false);
}

export function CertificatePolicyFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificatePolicy {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'keySpec': JsonWebKeySpecFromJSON(json['keySpec']),
        'ext': json['ext'],
        'expiryTime': json['expiryTime'],
        'alg': JsonWebSignatureAlgorithmFromJSON(json['alg']),
        'issuerIdentifier': json['issuerIdentifier'],
        'subject': CertificateSubjectFromJSON(json['subject']),
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : SubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'flags': ((json['flags'] as Array<any>).map(CertificateFlagFromJSON)),
    };
}

export function CertificatePolicyToJSON(value?: CertificatePolicy | null): any {
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
        'keySpec': JsonWebKeySpecToJSON(value.keySpec),
        'ext': value.ext,
        'expiryTime': value.expiryTime,
        'alg': JsonWebSignatureAlgorithmToJSON(value.alg),
        'issuerIdentifier': value.issuerIdentifier,
        'subject': CertificateSubjectToJSON(value.subject),
        'subjectAlternativeNames': SubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'flags': ((value.flags as Array<any>).map(CertificateFlagToJSON)),
    };
}

