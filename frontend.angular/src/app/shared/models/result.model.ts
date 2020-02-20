export class Result<T> {
  public success: boolean;
  public message: string;
  public value: T;
}

export class ListResult<T> extends Result<T> {
  public count:  number;
}
