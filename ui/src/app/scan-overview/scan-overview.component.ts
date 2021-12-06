import { Component, OnInit, Input } from '@angular/core';
import { FileData, Error } from '../search-binary.service';

@Component({
  selector: 'app-scan-overview',
  templateUrl: './scan-overview.component.html',
  styleUrls: ['./scan-overview.component.scss']
})

export class ScanOverviewComponent implements OnInit {

  @Input() data: any;
  @Input() checkName: any;
  errors: Error[] = [{repository: 'gus', error: 'you suck'}];
  checks: any[] = [];
  flattenedData: any[] = [];

  constructor() { }

  ngOnInit(): void {
    let result: any[] = [];
    this.data.forEach((f1: any) => {
      f1["checks"].forEach((check: any) => {
          let sahne = this.mapCheck(check)
          Object.keys(f1).forEach(ku => {
              if (ku != "checks") {
                  sahne[ku] = f1[ku];
              }
          })
          result.push(sahne);
      })
    });
    console.log(result);
    this.flattenedData = Object.assign([], result);
  }

  mapCheck(gus: any): any {
      let newObj: any = {}
      Object.keys(gus).forEach((key1: any) => {
          if (key1 == "additionalInfo") {
              Object.keys(gus[key1]).forEach((key2: any) => {
                newObj[key2] = gus[key1][key2];
              });
          } else {
            newObj[key1] = gus[key1];
          }
      })
      return newObj;
  }

  isBinariesSearch() {
    return this.checkName.includes('SearchBinaries');
  }

  isUnicodeSearch() {
    return this.checkName.includes('SearchUnicode');
  }

  isCommitSearch() {
    return this.checkName.includes('CheckCommitAuthor');
  }

  isBigFileSearch() {
    return this.checkName.includes('SearchBigFiles');
  }

  getCheckName(checkNamePath: string): string {
    checkNamePath.slice(checkNamePath.lastIndexOf('.') + 1);
    return checkNamePath;
  }

  listErrors() {
    this.data.array.forEach((entry: FileData) => {
      if (entry.error && entry.repository && entry.error !== '') {
        this.errors.push({repository: entry.repository, error: entry.error})
      }
    });
  }

}
