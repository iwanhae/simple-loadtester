export interface ResponseResult {
    isFailed: boolean;
    startedAt?: Date;
    duration?: number;

    color?: string;
    status?: number;
}