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


/**
 * 
 * @export
 */
export const LinkedCertificateTemplateUsage = {
    LinkedCertificateTemplateUsageClientAuthorization: 'cliant-authorization',
    LinkedCertificateTemplateUsageMemberDelegatedEnrollment: 'member-delegated-enrollment'
} as const;
export type LinkedCertificateTemplateUsage = typeof LinkedCertificateTemplateUsage[keyof typeof LinkedCertificateTemplateUsage];


export function LinkedCertificateTemplateUsageFromJSON(json: any): LinkedCertificateTemplateUsage {
    return LinkedCertificateTemplateUsageFromJSONTyped(json, false);
}

export function LinkedCertificateTemplateUsageFromJSONTyped(json: any, ignoreDiscriminator: boolean): LinkedCertificateTemplateUsage {
    return json as LinkedCertificateTemplateUsage;
}

export function LinkedCertificateTemplateUsageToJSON(value?: LinkedCertificateTemplateUsage | null): any {
    return value as any;
}
