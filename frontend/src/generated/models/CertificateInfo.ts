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
import type { CertificateUsage } from './CertificateUsage';
import {
    CertificateUsageFromJSON,
    CertificateUsageFromJSONTyped,
    CertificateUsageToJSON,
} from './CertificateUsage';
import type { JwkProperties } from './JwkProperties';
import {
    JwkPropertiesFromJSON,
    JwkPropertiesFromJSONTyped,
    JwkPropertiesToJSON,
} from './JwkProperties';

/**
 * 
 * @export
 * @interface CertificateInfo
 */
export interface CertificateInfo {
    /**
     * 
     * @type {string}
     * @memberof CertificateInfo
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfo
     */
    locator: string;
    /**
     * Time when the resoruce was last updated
     * @type {Date}
     * @memberof CertificateInfo
     */
    updated?: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfo
     */
    updatedBy?: string;
    /**
     * Time when the deleted was deleted
     * @type {Date}
     * @memberof CertificateInfo
     */
    deleted?: Date;
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof CertificateInfo
     */
    metadata?: { [key: string]: any; };
    /**
     * X.509 certificate SHA-1 thumbprint
     * @type {string}
     * @memberof CertificateInfo
     */
    thumbprint: string;
    /**
     * Common name
     * @type {string}
     * @memberof CertificateInfo
     */
    subjectCommonName: string;
    /**
     * Expiration date of the certificate
     * @type {Date}
     * @memberof CertificateInfo
     */
    notAfter: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfo
     */
    template: string;
    /**
     * Whether the certificate has been issued
     * @type {boolean}
     * @memberof CertificateInfo
     */
    isIssued: boolean;
    /**
     * Expiration date of the certificate
     * @type {Date}
     * @memberof CertificateInfo
     */
    notBefore: Date;
    /**
     * 
     * @type {Array<CertificateUsage>}
     * @memberof CertificateInfo
     */
    usages: Array<CertificateUsage>;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfo
     */
    issuer: string;
    /**
     * 
     * @type {JwkProperties}
     * @memberof CertificateInfo
     */
    jwk: JwkProperties;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfo
     */
    pem?: string;
}

/**
 * Check if a given object implements the CertificateInfo interface.
 */
export function instanceOfCertificateInfo(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "locator" in value;
    isInstance = isInstance && "thumbprint" in value;
    isInstance = isInstance && "subjectCommonName" in value;
    isInstance = isInstance && "notAfter" in value;
    isInstance = isInstance && "template" in value;
    isInstance = isInstance && "isIssued" in value;
    isInstance = isInstance && "notBefore" in value;
    isInstance = isInstance && "usages" in value;
    isInstance = isInstance && "issuer" in value;
    isInstance = isInstance && "jwk" in value;

    return isInstance;
}

export function CertificateInfoFromJSON(json: any): CertificateInfo {
    return CertificateInfoFromJSONTyped(json, false);
}

export function CertificateInfoFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateInfo {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'locator': json['locator'],
        'updated': !exists(json, 'updated') ? undefined : (new Date(json['updated'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'metadata': !exists(json, 'metadata') ? undefined : json['metadata'],
        'thumbprint': json['thumbprint'],
        'subjectCommonName': json['subjectCommonName'],
        'notAfter': (new Date(json['notAfter'])),
        'template': json['template'],
        'isIssued': json['isIssued'],
        'notBefore': (new Date(json['notBefore'])),
        'usages': ((json['usages'] as Array<any>).map(CertificateUsageFromJSON)),
        'issuer': json['issuer'],
        'jwk': JwkPropertiesFromJSON(json['jwk']),
        'pem': !exists(json, 'pem') ? undefined : json['pem'],
    };
}

export function CertificateInfoToJSON(value?: CertificateInfo | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'locator': value.locator,
        'updated': value.updated === undefined ? undefined : (value.updated.toISOString()),
        'updatedBy': value.updatedBy,
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'metadata': value.metadata,
        'thumbprint': value.thumbprint,
        'subjectCommonName': value.subjectCommonName,
        'notAfter': (value.notAfter.toISOString()),
        'template': value.template,
        'isIssued': value.isIssued,
        'notBefore': (value.notBefore.toISOString()),
        'usages': ((value.usages as Array<any>).map(CertificateUsageToJSON)),
        'issuer': value.issuer,
        'jwk': JwkPropertiesToJSON(value.jwk),
        'pem': value.pem,
    };
}

