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
 * @interface CertPolicy
 */
export interface CertPolicy {
    /**
     * 
     * @type {string}
     * @memberof CertPolicy
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof CertPolicy
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof CertPolicy
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof CertPolicy
     */
    updatedBy?: string;
    /**
     * 
     * @type {string}
     * @memberof CertPolicy
     */
    displayName: string;
    /**
     * 
     * @type {NamespaceKind}
     * @memberof CertPolicy
     * @deprecated
     */
    issuerNamespaceKind: NamespaceKind;
    /**
     * 
     * @type {string}
     * @memberof CertPolicy
     */
    issuerNamespaceIdentifier: string;
    /**
     * 
     * @type {SigningKeySpec}
     * @memberof CertPolicy
     */
    keySpec: SigningKeySpec;
    /**
     * 
     * @type {boolean}
     * @memberof CertPolicy
     */
    keyExportable: boolean;
    /**
     * 
     * @type {string}
     * @memberof CertPolicy
     */
    expiryTime: string;
    /**
     * 
     * @type {LifetimeAction}
     * @memberof CertPolicy
     */
    lifetimeAction?: LifetimeAction;
    /**
     * 
     * @type {CertificateSubject}
     * @memberof CertPolicy
     */
    subject: CertificateSubject;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof CertPolicy
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateFlag>}
     * @memberof CertPolicy
     */
    flags: Array<CertificateFlag>;
    /**
     * 
     * @type {string}
     * @memberof CertPolicy
     */
    version: string;
}

/**
 * Check if a given object implements the CertPolicy interface.
 */
export function instanceOfCertPolicy(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "displayName" in value;
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

export function CertPolicyFromJSON(json: any): CertPolicy {
    return CertPolicyFromJSONTyped(json, false);
}

export function CertPolicyFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertPolicy {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'displayName': json['displayName'],
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

export function CertPolicyToJSON(value?: CertPolicy | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
        'displayName': value.displayName,
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

