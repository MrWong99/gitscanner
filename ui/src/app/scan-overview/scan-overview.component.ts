import { Component, OnChanges, OnInit, Input, OnDestroy } from '@angular/core';
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
  //cols: any[] = [];
  //_selectedColumns: any[] = [];

  constructor(
    private messageService: MessageService,
    private searchBinaryService: SearchBinaryService) {
  }

  /*ngOnInit() {
    this.cols = [
      { field: 'authorEmail', header: 'Author Email' },
      { field: 'authorName', header: 'Author Name' },
      { field: 'character', header: 'Character' },
      { field: 'commiterEmail', header: 'Commiter Email' },
      { field: 'commitMessage', header: 'Commit Message' },
      { field: 'commiterName', header: 'Commiter Name' },
      { field: 'commitSize', header: 'Commit Size' },
      { field: 'filemode', header: 'File Mode' },
      { field: 'filesize', header: 'File Size' },
      { field: 'numberOfParents', header: 'Number Of Parents' }
    ];

    this.selectedColumns = this.cols;
  }*/

  ngOnChanges() {
    // flatten the nested data object retrieved from the parent/backend
    let result: any[] = [];
    this.data.forEach((repoCheck: any) => {
      if (repoCheck['checks']) {
        repoCheck['checks'].forEach((check: any) => {
            let flattenedCheck = this.mapCheck(check)
            Object.keys(repoCheck).forEach(k => {
                if (k !== 'checks' && k !== 'id') {
                  flattenedCheck[k] = repoCheck[k];
                }
            });
            /*if (flattenedCheck['checkName'] === 'SearchBinaries') {
              this.selectedColumns.push(this.cols[7]);
              this.selectedColumns.push(this.cols[8]);
            }*/
            result.push(flattenedCheck);
        });
      }
    });
    this.flattenedData = Object.assign([], result);
    console.log(this.flattenedData);
    this.listErrors();
  }

  /*@Input() get selectedColumns(): any[] {
    return this._selectedColumns;
  }

  set selectedColumns(val: any[]) {
      //restore original order
      this._selectedColumns = this.cols.filter(col => val.includes(col));
  }*/

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
   * Check if  selected check has any results to display.
   * @param checkName the check name
   * @returns true/false
   */
  hasResultData(checkName: string): number {
    let hasData: number = 0;
    this.flattenedData.forEach(data => {
      if (data.checkName === checkName) {
        hasData++;
      }
    });
    return hasData;
  }

  /**
   * Check if is a binary search check.
   * @returns true/false
   */
  isBinariesSearch(): boolean {
    return this.checkName.includes('SearchBinaries') && this.hasResultData('SearchBinaries') > 0;
  }

   /**
   * Check if is unicode check.
   * @returns true/false
   */
  isUnicodeSearch(): boolean {
    return this.checkName.includes('SearchIllegalUnicodeCharacters') && this.hasResultData('SearchIllegalUnicodeCharacters') > 0;
  }

   /**
   * Check if is commit author check.
   * @returns true/false
   */
  isCommitSearch(): boolean {
    return this.checkName.includes('CheckCommitMetaInformation')  && this.hasResultData('CheckCommitMetaInformation') > 0;
  }

   /**
   * Check if is big files check.
   * @returns true/false
   */
  isBigFileSearch(): boolean {
    return this.checkName.includes('SearchBigFiles') && this.hasResultData('SearchBigFiles') > 0;
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
