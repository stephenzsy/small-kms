/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
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
import type { LifetimeAction } from './LifetimeAction';
import {
    LifetimeActionFromJSON,
    LifetimeActionFromJSONTyped,
    LifetimeActionToJSON,
} from './LifetimeAction';
import type { NamespaceKind } from './NamespaceKind';
import {
    NamespaceKindFromJSON,
    NamespaceKindFromJSONTyped,
    NamespaceKindToJSON,
} from './NamespaceKind';
import type { SigningKeySpec } from './SigningKeySpec';
import {
    SigningKeySpecFromJSON,
    SigningKeySpecFromJSONTyped,
    SigningKeySpecToJSON,
} from './SigningKeySpec';
import type { SubjectAlternativeNames } from './SubjectAlternativeNames';
import {
    SubjectAlternativeNamesFromJSON,
    SubjectAlternativeNamesFromJSONTyped,
    SubjectAlternativeNamesToJSON,
} from './SubjectAlternativeNames';

/**
 * 
 * @export
 * @interface CertPolicyFields
 */
export interface CertPolicyFields {
    /**
     * 
     * @type {NamespaceKind}
     * @memberof CertPolicyFields
     */
    issuerNamespaceKind: NamespaceKind;
    /**
     * 
     * @type {string}
     * @memberof CertPolicyFields
     */
    issuerNamespaceIdentifier: string;
    /**
     * 
     * @type {SigningKeySpec}
     * @memberof CertPolicyFields
     */
    keySpec: SigningKeySpec;
    /**
     * 
     * @type {boolean}
     * @memberof CertPolicyFields
     */
    keyExportable: boolean;
    /**
     * 
     * @type {string}
     * @memberof CertPolicyFields
     */
    expiryTime: string;
    /**
     * 
     * @type {LifetimeAction}
     * @memberof CertPolicyFields
     */
    lifetimeAction?: LifetimeAction;
    /**
     * 
     * @type {CertificateSubject}
     * @memberof CertPolicyFields
     */
    subject: CertificateSubject;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof CertPolicyFields
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateFlag>}
     * @memberof CertPolicyFields
     */
    flags: Array<CertificateFlag>;
    /**
     * 
     * @type {string}
     * @memberof CertPolicyFields
     */
    version: string;
}

/**
 * Check if a given object implements the CertPolicyFields interface.
 */
export function instanceOfCertPolicyFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "issuerNamespaceKind" in value;
    isInstance = isInstance && "issuerNamespaceIdentifier" in value;
    isInstance = isInstance && "keySpec" in value;
    isInstance = isInstance && "keyExportable" in value;
    isInstance = isInstance && "expiryTime" in value;
    isInstance = isInstance && "subject" in value;
    isInstance = isInstance && "flags" in value;
    isInstance = isInstance && "version" in value;

    return isInstance;
}

export function CertPolicyFieldsFromJSON(json: any): CertPolicyFields {
    return CertPolicyFieldsFromJSONTyped(json, false);
}

export function CertPolicyFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertPolicyFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'issuerNamespaceKind': NamespaceKindFromJSON(json['issuerNamespaceKind']),
        'issuerNamespaceIdentifier': json['issuerNamespaceIdentifier'],
        'keySpec': SigningKeySpecFromJSON(json['keySpec']),
        'keyExportable': json['keyExportable'],
        'expiryTime': json['expiryTime'],
        'lifetimeAction': !exists(json, 'lifetimeAction') ? undefined : LifetimeActionFromJSON(json['lifetimeAction']),
        'subject': CertificateSubjectFromJSON(json['subject']),
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : SubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'flags': ((json['flags'] as Array<any>).map(CertificateFlagFromJSON)),
        'version': json['version'],
    };
}

export function CertPolicyFieldsToJSON(value?: CertPolicyFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'issuerNamespaceKind': NamespaceKindToJSON(value.issuerNamespaceKind),
        'issuerNamespaceIdentifier': value.issuerNamespaceIdentifier,
        'keySpec': SigningKeySpecToJSON(value.keySpec),
        'keyExportable': value.keyExportable,
        'expiryTime': value.expiryTime,
        'lifetimeAction': LifetimeActionToJSON(value.lifetimeAction),
        'subject': CertificateSubjectToJSON(value.subject),
        'subjectAlternativeNames': SubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'flags': ((value.flags as Array<any>).map(CertificateFlagToJSON)),
        'version': value.version,
    };
}

