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
    // flatten the nested data object retrieved from the parent/backend
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

  /**
   * Flatten a given check object.
   * @param check the check
   * @returns a flattened object
   */
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

  /**
   * Check if is a binary search check.
   * @returns true/false
   */
  isBinariesSearch(): boolean {
    return this.checkName.includes('SearchBinaries');
  }

   /**
   * Check if is unicode check.
   * @returns true/false
   */
  isUnicodeSearch(): boolean {
    return this.checkName.includes('SearchUnicode');
  }

   /**
   * Check if is commit author check.
   * @returns true/false
   */
  isCommitSearch(): boolean {
    return this.checkName.includes('CheckCommitAuthor');
  }

   /**
   * Check if is big files check.
   * @returns true/false
   */
  isBigFileSearch(): boolean {
    return this.checkName.includes('SearchBigFiles');
  }

  /**
   * List any errors that occured during the scan.
   */
  listErrors() {
    this.errors = [];
    this.data.forEach((entry: FileData) => {
      if (entry.error && entry.repository && entry.error !== '') {
        this.errors.push({repository: entry.repository, error: entry.error})
      }
    });
  }

  /**
   * Update the acknowledged status for the selected check.
   * @param id the check id
   * @param isAcknowledged true/false
   */
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
