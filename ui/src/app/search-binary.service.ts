import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface CheckNames {
  key: string,
  name: string
}

export interface Error {
  repository: string,
  error: string
}

export interface FileData {
  date: string,
  repository: string,
  error: string,
  checks: Checks[]
}

interface Checks {
  origin: string,
  branch: string,
  checkName: string,
  acknowledged: boolean,
  additionalInfo: AdditionalInfo
}

interface AdditionalInfo {
  authorName?: string,
  authorEmail?: string,
  character?: string,
  commitMessage?: string,
  commiterEmail?: string,
  commiterName?: string,
  filemode?: string,
  filesize?: string,
  numberOfParents?: number
}

@Injectable({
  providedIn: 'root'
})
export class SearchBinaryService {
  baseUrl: string = '/api/v1/';

  constructor(private httpClient: HttpClient) { }

  /**
   * Get the file data.
   * @param path path to scan
   * @param mode output format (path, type, size, full)
   * @returns array with the requested file data
   */
  getFileData(path: string, checkNames: string[]): Observable<FileData[]> {
    let body = {path: path, checkNames: checkNames};
    let headers = new HttpHeaders();
    headers.set('Content-Type', 'application/json');
    // binary files, illegal Unicode chars, commit email
    // binary files: github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries
    // illegal unicode chars:
    return this.httpClient.post<FileData[]>(this.baseUrl + 'checkRepos', body, {headers: headers});
  }

  updateAcknowledgedStatus(id: number, isAcknowledged: boolean): Observable<any> {
    let body = {acknowledged: isAcknowledged};
    let headers = new HttpHeaders();
    headers.set('Content-Type', 'application/json');
    return this.httpClient.put<any>(this.baseUrl + 'acknowledged/' + id, body, {headers: headers});
  }

  getStoredChecks(from: Date, to: Date, checkNames: string[]): Observable<FileData[]> {
    const fromDate = from.getMilliseconds();
    const toDate = to.getMilliseconds();
    const checks = checkNames.join(',');
    let params = new HttpParams();
    params.append('from', fromDate.toString());
    params.append('to', toDate.toString());
    params.append('checkNames', checks);
    return this.httpClient.get<FileData[]>(this.baseUrl + 'checks', {params: params});
  }
}
