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
import type { CertificateStatus } from './CertificateStatus';
import {
    CertificateStatusFromJSON,
    CertificateStatusFromJSONTyped,
    CertificateStatusToJSON,
} from './CertificateStatus';
import type { JsonWebSignatureKey } from './JsonWebSignatureKey';
import {
    JsonWebSignatureKeyFromJSON,
    JsonWebSignatureKeyFromJSONTyped,
    JsonWebSignatureKeyToJSON,
} from './JsonWebSignatureKey';
import type { SubjectAlternativeNames } from './SubjectAlternativeNames';
import {
    SubjectAlternativeNamesFromJSON,
    SubjectAlternativeNamesFromJSONTyped,
    SubjectAlternativeNamesToJSON,
} from './SubjectAlternativeNames';

/**
 * 
 * @export
 * @interface Certificate
 */
export interface Certificate {
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof Certificate
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    updatedBy: string;
    /**
     * 
     * @type {Date}
     * @memberof Certificate
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    displayName?: string;
    /**
     * Hex encoded certificate thumbprint
     * @type {string}
     * @memberof Certificate
     */
    thumbprint: string;
    /**
     * 
     * @type {CertificateStatus}
     * @memberof Certificate
     */
    status?: CertificateStatus;
    /**
     * 
     * @type {number}
     * @memberof Certificate
     */
    iat?: number;
    /**
     * 
     * @type {number}
     * @memberof Certificate
     */
    exp: number;
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    policyIdentifier: string;
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    issuerIdentifier: string;
    /**
     * 
     * @type {number}
     * @memberof Certificate
     */
    nbf: number;
    /**
     * 
     * @type {JsonWebSignatureKey}
     * @memberof Certificate
     */
    jwk?: JsonWebSignatureKey;
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    subject: string;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof Certificate
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateFlag>}
     * @memberof Certificate
     */
    flags?: Array<CertificateFlag>;
    /**
     * Key Vault certificate ID
     * @type {string}
     * @memberof Certificate
     */
    cid?: string;
    /**
     * Key Vault Secret ID
     * @type {string}
     * @memberof Certificate
     */
    sid?: string;
}

/**
 * Check if a given object implements the Certificate interface.
 */
export function instanceOfCertificate(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "thumbprint" in value;
    isInstance = isInstance && "exp" in value;
    isInstance = isInstance && "policyIdentifier" in value;
    isInstance = isInstance && "issuerIdentifier" in value;
    isInstance = isInstance && "nbf" in value;
    isInstance = isInstance && "subject" in value;

    return isInstance;
}

export function CertificateFromJSON(json: any): Certificate {
    return CertificateFromJSONTyped(json, false);
}

export function CertificateFromJSONTyped(json: any, ignoreDiscriminator: boolean): Certificate {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'thumbprint': json['thumbprint'],
        'status': !exists(json, 'status') ? undefined : CertificateStatusFromJSON(json['status']),
        'iat': !exists(json, 'iat') ? undefined : json['iat'],
        'exp': json['exp'],
        'policyIdentifier': json['policyIdentifier'],
        'issuerIdentifier': json['issuerIdentifier'],
        'nbf': json['nbf'],
        'jwk': !exists(json, 'jwk') ? undefined : JsonWebSignatureKeyFromJSON(json['jwk']),
        'subject': json['subject'],
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : SubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'flags': !exists(json, 'flags') ? undefined : ((json['flags'] as Array<any>).map(CertificateFlagFromJSON)),
        'cid': !exists(json, 'cid') ? undefined : json['cid'],
        'sid': !exists(json, 'sid') ? undefined : json['sid'],
    };
}

export function CertificateToJSON(value?: Certificate | null): any {
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
        'thumbprint': value.thumbprint,
        'status': CertificateStatusToJSON(value.status),
        'iat': value.iat,
        'exp': value.exp,
        'policyIdentifier': value.policyIdentifier,
        'issuerIdentifier': value.issuerIdentifier,
        'nbf': value.nbf,
        'jwk': JsonWebSignatureKeyToJSON(value.jwk),
        'subject': value.subject,
        'subjectAlternativeNames': SubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'flags': value.flags === undefined ? undefined : ((value.flags as Array<any>).map(CertificateFlagToJSON)),
        'cid': value.cid,
        'sid': value.sid,
    };
}

