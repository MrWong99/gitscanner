import { Component, OnChanges, Input, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { MessageService } from 'primeng/api';
import { SearchBinaryService, FileData, Error } from '../search-binary.service';

@Component({
  selector: 'app-scan-overview',
  templateUrl: './scan-overview.component.html',
  styleUrls: ['./scan-overview.component.scss']
})

export class ScanOverviewComponent implements OnChanges, OnDestroy {

  @Input() data: any;
  @Input() checkName: any;
  errors: Error[] = [];
  checks: any[] = [];
  flattenedData: any[] = [];
  updateAcknowledgedStatusSub: Subscription | undefined;

  constructor(
    private messageService: MessageService,
    private searchBinaryService: SearchBinaryService) {
  }

  ngOnChanges() {
    let result: any[] = [];
    this.data.forEach((repoCheck: any) => {
      if (repoCheck['checks']) {
        repoCheck['checks'].forEach((check: any) => {
            let flattenedCheck = this.mapCheck(check)
            Object.keys(repoCheck).forEach(k => {
                if (k != 'checks' && k != 'id') {
                  flattenedCheck[k] = repoCheck[k];
                }
            })
            result.push(flattenedCheck);
        });
      }
    });
    this.flattenedData = Object.assign([], result);
    this.listErrors();
  }

  mapCheck(check: any): any {
      let flattenedObj: any = {}
      Object.keys(check).forEach((key1: any) => {
          if (key1 == 'additionalInfo') {
              Object.keys(check[key1]).forEach((key2: any) => {
                flattenedObj[key2] = check[key1][key2];
              });
          } else {
            flattenedObj[key1] = check[key1];
          }
      })
      return flattenedObj;
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
    this.errors = [];
    this.data.forEach((entry: FileData) => {
      if (entry.error && entry.repository && entry.error !== '') {
        this.errors.push({repository: entry.repository, error: entry.error})
      }
    });
  }

  updateAcknowledgedStatus(id: number, isAcknowledged: boolean) {
    this.updateAcknowledgedStatusSub = this.searchBinaryService.updateAcknowledgedStatus(id, isAcknowledged).subscribe(data => {
     },
     error => {
      this.messageService.add({key: 'acknowledgedError', severity:'error', summary:'Error updating acknowledged status', detail: error.statusText});
     });
  }

  ngOnDestroy() {
    this.updateAcknowledgedStatusSub?.unsubscribe();
  }

}
