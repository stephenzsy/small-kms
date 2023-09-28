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
import type { NamespaceTypeShortName } from './NamespaceTypeShortName';
import {
    NamespaceTypeShortNameFromJSON,
    NamespaceTypeShortNameFromJSONTyped,
    NamespaceTypeShortNameToJSON,
} from './NamespaceTypeShortName';

/**
 * 
 * @export
 * @interface CertificateIssuer
 */
export interface CertificateIssuer {
    /**
     * 
     * @type {NamespaceTypeShortName}
     * @memberof CertificateIssuer
     */
    namespaceType: NamespaceTypeShortName;
    /**
     * 
     * @type {string}
     * @memberof CertificateIssuer
     */
    namespaceId: string;
    /**
     * if certificate ID is not specified, use template ID to find the latest certificate, use default value if not specified
     * @type {string}
     * @memberof CertificateIssuer
     */
    templateId?: string;
}

/**
 * Check if a given object implements the CertificateIssuer interface.
 */
export function instanceOfCertificateIssuer(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceType" in value;
    isInstance = isInstance && "namespaceId" in value;

    return isInstance;
}

export function CertificateIssuerFromJSON(json: any): CertificateIssuer {
    return CertificateIssuerFromJSONTyped(json, false);
}

export function CertificateIssuerFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateIssuer {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'namespaceType': NamespaceTypeShortNameFromJSON(json['namespaceType']),
        'namespaceId': json['namespaceId'],
        'templateId': !exists(json, 'templateId') ? undefined : json['templateId'],
    };
}

export function CertificateIssuerToJSON(value?: CertificateIssuer | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'namespaceType': NamespaceTypeShortNameToJSON(value.namespaceType),
        'namespaceId': value.namespaceId,
        'templateId': value.templateId,
    };
}

