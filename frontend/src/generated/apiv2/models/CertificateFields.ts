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
import type { CertificatePendingAcme } from './CertificatePendingAcme';
import {
    CertificatePendingAcmeFromJSON,
    CertificatePendingAcmeFromJSONTyped,
    CertificatePendingAcmeToJSON,
} from './CertificatePendingAcme';
import type { JsonWebKey } from './JsonWebKey';
import {
    JsonWebKeyFromJSON,
    JsonWebKeyFromJSONTyped,
    JsonWebKeyToJSON,
} from './JsonWebKey';
import type { SubjectAlternativeNames } from './SubjectAlternativeNames';
import {
    SubjectAlternativeNamesFromJSON,
    SubjectAlternativeNamesFromJSONTyped,
    SubjectAlternativeNamesToJSON,
} from './SubjectAlternativeNames';

/**
 * 
 * @export
 * @interface CertificateFields
 */
export interface CertificateFields {
    /**
     * 
     * @type {string}
     * @memberof CertificateFields
     */
    identififier: string;
    /**
     * 
     * @type {string}
     * @memberof CertificateFields
     */
    issuerIdentifier: string;
    /**
     * 
     * @type {number}
     * @memberof CertificateFields
     */
    nbf: number;
    /**
     * 
     * @type {JsonWebKey}
     * @memberof CertificateFields
     */
    jwk?: JsonWebKey;
    /**
     * 
     * @type {string}
     * @memberof CertificateFields
     */
    subject: string;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof CertificateFields
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateFlag>}
     * @memberof CertificateFields
     */
    flags?: Array<CertificateFlag>;
    /**
     * Key Vault certificate ID
     * @type {string}
     * @memberof CertificateFields
     */
    cid?: string;
    /**
     * Key Vault Secret ID
     * @type {string}
     * @memberof CertificateFields
     */
    sid?: string;
    /**
     * 
     * @type {JsonWebKey}
     * @memberof CertificateFields
     */
    oneTimePkcs12Key?: JsonWebKey;
    /**
     * 
     * @type {CertificatePendingAcme}
     * @memberof CertificateFields
     */
    pendingAcme?: CertificatePendingAcme;
}

/**
 * Check if a given object implements the CertificateFields interface.
 */
export function instanceOfCertificateFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "identififier" in value;
    isInstance = isInstance && "issuerIdentifier" in value;
    isInstance = isInstance && "nbf" in value;
    isInstance = isInstance && "subject" in value;

    return isInstance;
}

export function CertificateFieldsFromJSON(json: any): CertificateFields {
    return CertificateFieldsFromJSONTyped(json, false);
}

export function CertificateFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'identififier': json['identififier'],
        'issuerIdentifier': json['issuerIdentifier'],
        'nbf': json['nbf'],
        'jwk': !exists(json, 'jwk') ? undefined : JsonWebKeyFromJSON(json['jwk']),
        'subject': json['subject'],
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : SubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'flags': !exists(json, 'flags') ? undefined : ((json['flags'] as Array<any>).map(CertificateFlagFromJSON)),
        'cid': !exists(json, 'cid') ? undefined : json['cid'],
        'sid': !exists(json, 'sid') ? undefined : json['sid'],
        'oneTimePkcs12Key': !exists(json, 'oneTimePkcs12Key') ? undefined : JsonWebKeyFromJSON(json['oneTimePkcs12Key']),
        'pendingAcme': !exists(json, 'pendingAcme') ? undefined : CertificatePendingAcmeFromJSON(json['pendingAcme']),
    };
}

export function CertificateFieldsToJSON(value?: CertificateFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'identififier': value.identififier,
        'issuerIdentifier': value.issuerIdentifier,
        'nbf': value.nbf,
        'jwk': JsonWebKeyToJSON(value.jwk),
        'subject': value.subject,
        'subjectAlternativeNames': SubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'flags': value.flags === undefined ? undefined : ((value.flags as Array<any>).map(CertificateFlagToJSON)),
        'cid': value.cid,
        'sid': value.sid,
        'oneTimePkcs12Key': JsonWebKeyToJSON(value.oneTimePkcs12Key),
        'pendingAcme': CertificatePendingAcmeToJSON(value.pendingAcme),
    };
}

