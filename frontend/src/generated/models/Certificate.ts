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
import type { CertificateUsage } from './CertificateUsage';
import {
    CertificateUsageFromJSON,
    CertificateUsageFromJSONTyped,
    CertificateUsageToJSON,
} from './CertificateUsage';

/**
 * 
 * @export
 * @interface Certificate
 */
export interface Certificate {
    /**
     * Unique ID of the namespace
     * @type {string}
     * @memberof Certificate
     */
    namespaceId: string;
    /**
     * 
     * @type {string}
     * @memberof Certificate
     */
    id: string;
    /**
     * Unique ID of the user who created the policy
     * @type {string}
     * @memberof Certificate
     */
    updatedBy: string;
    /**
     * Time when the policy was last updated
     * @type {Date}
     * @memberof Certificate
     */
    updated: Date;
    /**
     * Time when the policy was deleted
     * @type {Date}
     * @memberof Certificate
     */
    deleted?: Date;
    /**
     * Name of the certificate, also the common name (CN) in the subject of the certificate
     * @type {string}
     * @memberof Certificate
     */
    name: string;
    /**
     * 
     * @type {CertificateUsage}
     * @memberof Certificate
     */
    usage: CertificateUsage;
    /**
     * Expiration date of the certificate
     * @type {Date}
     * @memberof Certificate
     */
    notAfter: Date;
    /**
     * Issuer namespace ID
     * @type {string}
     * @memberof Certificate
     */
    issuerNamespace: string;
    /**
     * Issuer certificate ID
     * @type {string}
     * @memberof Certificate
     */
    issuer: string;
    /**
     * Unique ID of the user who created the certificate
     * @type {string}
     * @memberof Certificate
     */
    createdBy: string;
    /**
     * PEM encoded X.509 certificate
     * @type {string}
     * @memberof Certificate
     */
    x509pem?: string;
}

/**
 * Check if a given object implements the Certificate interface.
 */
export function instanceOfCertificate(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceId" in value;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "usage" in value;
    isInstance = isInstance && "notAfter" in value;
    isInstance = isInstance && "issuerNamespace" in value;
    isInstance = isInstance && "issuer" in value;
    isInstance = isInstance && "createdBy" in value;

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
        
        'namespaceId': json['namespaceId'],
        'id': json['id'],
        'updatedBy': json['updatedBy'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'name': json['name'],
        'usage': CertificateUsageFromJSON(json['usage']),
        'notAfter': (new Date(json['notAfter'])),
        'issuerNamespace': json['issuerNamespace'],
        'issuer': json['issuer'],
        'createdBy': json['createdBy'],
        'x509pem': !exists(json, 'x509pem') ? undefined : json['x509pem'],
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
        
        'namespaceId': value.namespaceId,
        'id': value.id,
        'updatedBy': value.updatedBy,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'name': value.name,
        'usage': CertificateUsageToJSON(value.usage),
        'notAfter': (value.notAfter.toISOString()),
        'issuerNamespace': value.issuerNamespace,
        'issuer': value.issuer,
        'createdBy': value.createdBy,
        'x509pem': value.x509pem,
    };
}
